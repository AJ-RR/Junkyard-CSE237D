package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1" //Use `go get` to install packages
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type JobResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
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
	/*
	*  SUBMIT Request Handler
	 */
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {

		//Check if request is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20) // 10MB max
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		//Read metadata of request (Student name, assignment name, etc.)
		student := strings.ToLower(r.FormValue("name"))
		assignment := strings.ToLower(r.FormValue("image"))
		if student == "" || assignment == "" {
			http.Error(w, "Missing 'name' or 'image' field", http.StatusBadRequest)
			return
		}
		name := fmt.Sprintf("%s-%s-%d", student, assignment, time.Now().Unix())
		startTime := time.Now() 

		//Read file from form into buffer
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

		job_ttl := int32(120) //How long to keep job alive after completion (120 seconds)
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
								Name: "runner",
								// Image: "arunanthivi/job-grader:python", //Python image for container
								Image:           "rsankar12/opencl_cse160", //Rishab's OpenCL image for container
								ImagePullPolicy: corev1.PullAlways,
								Command: []string{
  "sh", "-c",
  "unzip /scripts/archive.zip -d $HOME >/dev/null 2>&1 && " +
  "cd $HOME/CSE160Assignment2/PA2 && " +                 // adjust if needed
  "make -s run > /tmp/out 2>&1; RC=$?; " +
  // if make failed, still emit JSON Gradescope can parse
  "if [ $RC -ne 0 ]; then echo '{\"score\":0}' ; else cat /tmp/out ; fi",
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
		_, err = jobClient.Create(context.TODO(), job, meta.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create Job: %v", err), http.StatusInternalServerError)
			return
		}
		defer clientset.CoreV1().ConfigMaps("default").Delete(context.TODO(), configMapName, meta.DeleteOptions{})
		defer jobClient.Delete(context.TODO(), name, meta.DeleteOptions{})

		//Poll for Job completion
		for {
			job, err := jobClient.Get(context.TODO(), name, meta.GetOptions{})
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get job status: %v", err), http.StatusInternalServerError)
				return
			}

			if job.Status.Succeeded > 0 {
				break
			} else if job.Status.Failed > 0 {
				http.Error(w, "Job failed", http.StatusInternalServerError)
				return
			}

			time.Sleep(2 * time.Second)
		}

		// Get the pod name associated with the job
		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), meta.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", name),
		})
		if err != nil || len(pods.Items) == 0 {
			http.Error(w, "Failed to list pods for job", http.StatusInternalServerError)
			return
		}

		podName := pods.Items[0].Name

		// Fetch logs from the pod
		logReq := clientset.CoreV1().Pods("default").GetLogs(podName, &corev1.PodLogOptions{})
		logStream, err := logReq.Stream(context.TODO())
		if err != nil {
			http.Error(w, "Failed to stream pod logs", http.StatusInternalServerError)
			return
		}
		defer logStream.Close()

		logs, err := io.ReadAll(logStream)
		if err != nil {
			http.Error(w, "Failed to read pod logs", http.StatusInternalServerError)
			return
		}

		completion := time.Now()
		updateLatency(startTime, completion)
		w.Header().Set("Content-Type", "application/json")
		w.Write(logs)
		
	})

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
