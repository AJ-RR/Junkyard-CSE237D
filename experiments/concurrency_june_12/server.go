package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp" // Import the regexp package
	"strings"
	"sync" // For protecting the jobStore map
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// JobResponse is sent back to the client immediately after job creation.
type JobResponse struct {
	Status string `json:"status"`
	JobID  string `json:"job_id,omitempty"` // Add JobID for client to poll
	Error  string `json:"error,omitempty"`
}

// JobStatusPayload is sent back to the client when polling for status.
type JobStatusPayload struct {
	Status  string `json:"status"`            // "pending", "succeeded", "failed"
	Results string `json:"results,omitempty"` // Logs from the job
	Error   string `json:"error,omitempty"`   // Error message if job failed or logs couldn't be fetched
}

// JobInternalState holds the internal state of a job managed by this server.
type JobInternalState struct {
	Status  string    // "pending", "succeeded", "failed"
	Results []byte    // Raw logs from the job
	Error   error     // Go error object if any issue occurred
}

// jobStore is an in-memory map to store the state of all active jobs.
// IMPORTANT: For production, this should be replaced with a persistent store
// like Redis or a database, as map data will be lost on server restart.
var jobStore = make(map[string]*JobInternalState)
var jobStoreMutex sync.Mutex // Mutex to protect jobStore from concurrent access

// sanitizeK8sName converts a string to be RFC 1123 compliant (lowercase alphanumeric, '-', '.', and starts/ends with alphanumeric).
func sanitizeK8sName(s string) string {
    // Convert to lowercase
    s = strings.ToLower(s)
    // Replace any characters not allowed by RFC 1123 (alphanumeric, -, .) with a hyphen
    reg := regexp.MustCompile("[^a-z0-9.-]+")
    s = reg.ReplaceAllString(s, "-")
    // Trim leading/trailing hyphens/dots
    s = strings.Trim(s, "-.")
    // Replace multiple hyphens with a single hyphen
    s = strings.ReplaceAll(s, "--", "-")
    return s
}


func main() {
	// Load kubeconfig from default or env var
	var config *rest.Config
	var err error
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig != "" {
		fmt.Println("Using kubeconfig:", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Failed to load kubeconfig: %v", err)
		}
	} else {
		fmt.Println("Using in-cluster config")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to load in-cluster config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Submit Request Handler (`/submit`)
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		// Check if request is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Read metadata of request (Student name, assignment name, etc.)
		student := r.FormValue("name")
		assignment := r.FormValue("image")
		if student == "" || assignment == "" {
			http.Error(w, "Missing 'name' or 'image' field", http.StatusBadRequest)
			return
		}

		// --- FIX START: Sanitize student and assignment names ---
		student = sanitizeK8sName(student)
		assignment = sanitizeK8sName(assignment)
		// --- FIX END ---
		
		name := fmt.Sprintf("%s-%s-%d", student, assignment, time.Now().Unix())
		// startTime := time.Now() // Keep for potential latency tracking later

		// Read file from form into buffer
		file, _, err := r.FormFile("script")
		if err != nil {
			http.Error(w, "Missing 'script' file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		zipData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read script file", http.StatusInternalServerError)
			return
		}

		// Create a ConfigMap to hold the script
		configMapName := "script-cm-" + name
		_, err = clientset.CoreV1().ConfigMaps("default").Create(context.TODO(), &corev1.ConfigMap{
			ObjectMeta: meta.ObjectMeta{
				Name: configMapName,
			},
			BinaryData: map[string][]byte{
				"archive.zip": zipData,
			},
		}, meta.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create ConfigMap: %v", err), http.StatusInternalServerError)
			return
		}
		jobClient := clientset.BatchV1().Jobs("default")

		job_ttl := int32(120) // How long to keep job alive after completion (120 seconds)
		// Create the Job that runs the script
		job := &batchv1.Job{
			ObjectMeta: meta.ObjectMeta{
				Name: name,
			},
			Spec: batchv1.JobSpec{
				TTLSecondsAfterFinished: &job_ttl,
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						Volumes: []corev1.Volume{
							{
								Name: "script-volume",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: configMapName,
										},
									},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:            "runner",
								Image:           "rsankar12/opencl_cse160", // Rishab's OpenCL image for container
								ImagePullPolicy: corev1.PullAlways,
								Command: []string{
									"sh", "-c",
									// ① unzip silently
									"unzip /scripts/archive.zip -d $HOME >/dev/null 2>&1 && " +
										// ② locate the PA2 folder (first match) and cd into it, suppressing errors
										"PA2DIR=$(find $HOME -type d -name PA2 | head -n1) && " +
										"{ cd \"$PA2DIR\" 2>/dev/null || EXIT=1; } && " +
										// ③ run make, capture output and exit‐code only if cd succeeded
										"if [ \"$EXIT\" != \"1\" ]; then make -s run > /tmp/out 2>&1; EXIT=$?; fi; " +
										// ④ emit *only* JSON, then exit 0
										"if [ \"$EXIT\" != \"0\" ]; then echo '{\"score\":0}'; else cat /tmp/out; fi",
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "script-volume",
										MountPath: "/scripts",
									},
								},
							},
						},
					},
				},
			},
		}

		// Create the Kubernetes Job
		_, err = jobClient.Create(context.TODO(), job, meta.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create Job: %v", err), http.StatusInternalServerError)
			// Clean up configmap if job creation failed
			deleteErr := clientset.CoreV1().ConfigMaps("default").Delete(context.TODO(), configMapName, meta.DeleteOptions{})
			if deleteErr != nil {
				log.Printf("Warning: Failed to delete ConfigMap %s after job creation failure: %v", configMapName, deleteErr)
			}
			return
		}

		// Initialize job status in our store
		jobStoreMutex.Lock()
		jobStore[name] = &JobInternalState{
			Status: "pending",
		}
		jobStoreMutex.Unlock()

		// Start a goroutine to monitor the job and update its status
		go func(jobName string, cfgMapName string, clientset *kubernetes.Clientset) {
			log.Printf("Starting goroutine to monitor job %s", jobName)

			// Ensure ConfigMap and Job are eventually deleted after monitoring completes
			defer func() {
				log.Printf("Attempting to delete ConfigMap %s for job %s", cfgMapName, jobName)
				deleteErr := clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), cfgMapName, meta.DeleteOptions{})
				if deleteErr != nil {
					log.Printf("Error deleting ConfigMap %s: %v", cfgMapName, deleteErr)
				}

				log.Printf("Attempting to delete Job %s", jobName)
				deleteErr = clientset.BatchV1().Jobs("default").Delete(context.Background(), jobName, meta.DeleteOptions{})
				if deleteErr != nil {
					log.Printf("Error deleting Job %s: %v", jobName, deleteErr)
				}
			}()

			var finalStatus string
			var jobLogs []byte
			var jobError error

			// Poll for Job completion
			for {
				job, err := clientset.BatchV1().Jobs("default").Get(context.TODO(), jobName, meta.GetOptions{})
				if err != nil {
					// This might happen if the job was deleted unexpectedly or if there's a K8s API issue.
					jobError = fmt.Errorf("Failed to get job status for %s: %v", jobName, err)
					finalStatus = "failed"
					break
				}

				if job.Status.Succeeded > 0 {
					finalStatus = "succeeded"
					break
				} else if job.Status.Failed > 0 {
					finalStatus = "failed"
					jobError = fmt.Errorf("Job %s failed on Kubernetes", jobName)
					break
				}

				time.Sleep(2 * time.Second) // Poll every 2 seconds
			}

			// Fetch logs if the job succeeded or failed
			if finalStatus == "succeeded" || finalStatus == "failed" {
				pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), meta.ListOptions{
					LabelSelector: fmt.Sprintf("job-name=%s", jobName),
				})
				if err != nil || len(pods.Items) == 0 {
					jobError = fmt.Errorf("Failed to list pods for job %s: %v", jobName, err)
					// Keep finalStatus as it was, but add log fetching error
				} else {
					podName := pods.Items[0].Name // Assuming one pod per job
					logReq := clientset.CoreV1().Pods("default").GetLogs(podName, &corev1.PodLogOptions{})
					logStream, err := logReq.Stream(context.TODO())
					if err != nil {
						jobError = fmt.Errorf("Failed to stream pod logs for %s (pod %s): %v", jobName, podName, err)
						// Keep finalStatus as it was, but add log fetching error
					} else {
						defer logStream.Close()
						logs, err := io.ReadAll(logStream)
						if err != nil {
							jobError = fmt.Errorf("Failed to read pod logs for %s (pod %s): %v", jobName, podName, err)
							// Keep finalStatus as it was, but add log fetching error
						} else {
							jobLogs = logs
						}
					}
				}
			}

			// Update job store with final status and results
			jobStoreMutex.Lock()
			if js, ok := jobStore[jobName]; ok {
				js.Status = finalStatus
				js.Results = jobLogs
				js.Error = jobError
			}
			jobStoreMutex.Unlock()
			log.Printf("Job %s completed with status: %s", jobName, finalStatus)

			// updateLatency(startTime, time.Now()) // If you want to log latency on server side
		}(name, configMapName, clientset) // Pass needed variables to the goroutine

		// Respond to the client immediately after creating the job
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) // 202 Accepted means request accepted for asynchronous processing
		json.NewEncoder(w).Encode(JobResponse{
			Status: "Job created, please poll /status/" + name + " for results",
			JobID:  name,
		})
	})

	// Status Request Handler (`/status/{jobName}`)
	http.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract jobName from the URL path, e.g., /status/my-job-name
		jobName := strings.TrimPrefix(r.URL.Path, "/status/")
		if jobName == "" {
			http.Error(w, "Missing job ID in URL path, e.g., /status/my-job-name", http.StatusBadRequest)
			return
		}

		jobStoreMutex.Lock()
		jobState, found := jobStore[jobName]
		jobStoreMutex.Unlock()

		if !found {
			http.Error(w, "Job ID not found or has been cleaned up", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		responsePayload := JobStatusPayload{
			Status: jobState.Status,
		}

		// Only include results/error if the job is actually complete
		if jobState.Status == "succeeded" || jobState.Status == "failed" {
			responsePayload.Results = string(jobState.Results)
			if jobState.Error != nil {
				responsePayload.Error = jobState.Error.Error()
			}
		}

		json.NewEncoder(w).Encode(responsePayload)
	})

	// Root Path Handler (`/`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Server is running"})
	})

	log.Println("Server listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

// Dummy updateLatency function - define as needed or remove if not used
