## Experiments Overview

This directory contains various prototype experiments and demo scripts to evaluate Kubernetes job submission, grading infrastructure, and job orchestration across different scenarios.

### Google_demo
Demonstrates creating a Kubernetes Job on a remote (Google Cloud) cluster using a Python client script.

Key files:
- `demo.py`: Uses the Kubernetes Python client to load a kubeconfig file, parse a job YAML spec, and submit the job to a specified namespace.
- `test_job.yaml`: Defines a simple Kubernetes job that launches a long-running container (currently just `sleep infinity`) using a custom OpenCL-based image.
- Verifies Kubernetes job submission, YAML parsing, and cluster connectivity on remote GKE clusters.

### cse160_assignment_demo
Demonstrates automated evaluation of the CSE160 Assignment 2 using Kubernetes Jobs.

Key files:
- `demo.py`: Python script that submits a Kubernetes job to a specified namespace using a provided kubeconfig file. It programmatically updates job metadata and loads job specs from a YAML file.
- `test_job.yaml`: Kubernetes job definition that runs a container based on a custom `opencl-kube` image. It:
  - Navigates to a specific assignment directory (`CSE160Assignment2/PA2`)
  - Executes `make run` and logs the output to `output.log`
  - Prints timestamps before and after execution
- Tests a realistic evaluation flow for student submissions.


### gradescope_scripts
Contains helper scripts that simulate and support Gradescope-style autograding within Kubernetes workflows.

Key files:
- `run_autograder`: Bash script that automates submission handling and grading integration. It:
  - Extracts student metadata (name, assignment title) from `submission_metadata.json`
  - Zips the student's submission files
  - Sends the submission to a remote grading server via a `curl` POST request
  - Saves the grading output to `results.json`
- `setup.sh`: Installs required system packages (`ssh`, `curl`, `jq`, and `zip`) inside the container environment.

### multiple_jobs_metrics


### python_client
