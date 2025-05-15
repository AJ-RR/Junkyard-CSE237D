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
		name := student + "-" + assignment

		//Read file from form into buffer
		file, _, err := r.FormFile("script")
		if err != nil {
			http.Error(w, "Missing 'script' file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		scriptData, err := io.ReadAll(file)
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
			Data: map[string]string{
				"main.py": string(scriptData),
			},
		}, meta.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create ConfigMap: %v", err), http.StatusInternalServerError)
			return
		}

		job_ttl := int32(120) //How long to keep job alive after completion
		// Create the Job that runs the script
		job := &batchv1.Job{
			ObjectMeta: meta.ObjectMeta{
				Name: name,
			},
			Spec: batchv1.JobSpec{
				TTLSecondsAfterFinished: &job_ttl, // Complete job after 2 minutes (120 seconds)
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
								Name:    "runner",
								Image:   "python:3.9", //Python image for container
								Command: []string{"python", "/scripts/main.py"},
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

		_, err = clientset.BatchV1().Jobs("default").Create(context.TODO(), job, meta.CreateOptions{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create Job: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JobResponse{Status: "Job created"})
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
