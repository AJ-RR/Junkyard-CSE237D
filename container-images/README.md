## container-images Overview

This directory contains Docker images and related components used for running and grading jobs inside containers.

### job-server
Contains the source code and Kubernetes manifests for a REST-based job submission server:
- Launches jobs in a Kubernetes cluster based on user-submitted ZIP files.
- Creates `ConfigMaps`, `Jobs`, and fetches logs from resulting pods.
- Includes latency tracking across jobs (`metrics.go`) and job lifecycle handling (`server.go`).
- Comes with a `Dockerfile` and deployment configuration (`jobserver.yaml`).

Key files:
- `server.go`: HTTP server that handles submissions and spawns Kubernetes Jobs.
- `metrics.go`: Tracks execution time metrics across jobs.
- `jobserver.yaml`: Kubernetes ServiceAccount, Role, RoleBinding, Deployment, and Service definitions.
- `Dockerfile`, `go.mod`, `go.sum`: For building and running the server in a container.

**Refer to [`Documentation/Setup-JobServer-Service`](./Documentation/Setup-Jobserver-Service/README.md)** for further detailed instructions

### python-grader

A containerized Python-based autograder using `gradescope-utils`.

Key components:
- `Dockerfile`: Builds an image from `python:3.13-slim`, installs `gradescope-utils`, and copies test harness and test cases.
- `harness.py`: Discovers and runs `unittest` test cases using Gradescope's `JSONTestRunner`.
- `tests/`: Directory containing Python unit tests that define pass/fail logic for student submissions.

### python-hello-world
A simple test container image setup used for verifying that the job submission infrastructure is working end-to-end.

