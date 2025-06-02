# Sample Assignment Evaluation: CSE160 Assignment 2

To evaluate the performance of the Green Grader, we used CSE160 Assignment 2 as a sample assignment. The following structure and components were involved:

- [`CSE160Assignment2`](./CSE160Assignment2):  
  This folder contains the main files related to the assignment used in testing:
  - `PA2/`: The primary assignment folder with the `Makefile` and main source code files submitted by students.
  - `helper_lib/`: A directory containing supporting helper files required for the assignment’s core functionality.

- [`rsankar12/opencl_cse160`](https://hub.docker.com/r/rsankar12/opencl_cse160):  
  This is the Docker container used for running the assignment.  
  - It is a customized version of Professor Kastner’s original CSE160 POCL container.
  - The container includes a modified `Makefile` and pre-installed helper files to match the system’s environment.

**Refer to [`Documentation/Setup-OpenCL`](./Documentation/Setup-OpenCL/README.md)** for detailed instructions on:
- Setting up the POCL container within a Kubernetes environment
- Modifications made to the original `Makefile` to adapt it to the operating system used in our testbed
