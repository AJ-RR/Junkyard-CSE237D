# Green Grader Documentation

# Documentation Overview

This directory contains documentation for each of the steps required to reporduce the Automated Grading System via the Phone Cluster. There is also a short description and link to the corresponding direcotry based on the specific information and steps you are looking for.

| Folder | Description |
|--------|-------------|
| [`Build-Android-Kernel`](./Documentation/Build-Android-Kernel) | Instructions to configure and build custom Android kernels for the reused phones. |
| [`Build-Custom-Image-Jobs`](./Documentation/Build-Custom-Image-Jobs) | Guides on building container or image-based job artifacts to be deployed across the phone cluster. |
| [`Connect-Gradescope-to-Cluster`](./Documentation/Connect-Gradescope-to-Cluster) | Integration steps for enabling Gradescope to send and receive grading jobs from the cluster. |
| [`Connect-to-Cluster`](./Documentation/Connect-to-Cluster) | Details on establishing secure and reliable connections to the phone-based cluster. |
| [`Setup-Jobserver-Service`](./Documentation/Setup-Jobserver-Service) | Setup instructions for the job server that manages and distributes tasks to phones. |
| [`Setup-Kubernetes`](./Documentation/Setup-Kubernetes) | (Optional) Steps for setting up a lightweight Kubernetes cluster to orchestrate workloads. |
| [`Setup-OpenCL`](./Documentation/Setup-OpenCL) | Documentation for enabling OpenCL GPU support on Android devices for accelerated tasks. |
| [`Setup-Webservers`](./Documentation/Setup-Webservers) | Deployment of web servers for APIs, monitoring dashboards, or application frontends. |
