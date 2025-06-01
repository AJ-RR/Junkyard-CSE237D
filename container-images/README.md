## container-images Overview

This directory contains Docker images and related components used for running and grading jobs inside containers.

### 📁 job-server
Contains the source code and Kubernetes manifests for a REST-based job submission server:
- Launches jobs in a Kubernetes cluster based on user-submitted ZIP files.
- Creates `ConfigMaps`, `Jobs`, and fetches logs from resulting pods.
- Includes latency tracking across jobs (`metrics.go`) and job lifecycle handling (`server.go`).
- Comes with a `Dockerfile` and deployment configuration (`jobserver.yaml`).
**Refer to [`Documentation/Setup-JobServer-Service`](./Documentation/Setup-JobServer-Service/README.md)** for further detailed instructions

### 📁 python-grader
Likely contains a Dockerfile and scripts for a Python-based grading container:
- Used by the `job-server` to execute student-submitted Python assignments.
- May include a `Makefile` or similar mechanism for automated grading logic.

### 📁 python-hello-world
A simple test container image setup:
- Likely includes a basic `hello.py` script or similar.
- Used for verifying that the job submission infrastructure is working end-to-end.

### 📄 README.md
Documentation for building and understanding the container images in this directory.

