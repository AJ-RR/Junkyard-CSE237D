## Setting Up OpenCL
### opencl-cpu.yaml
```
apiVersion: v1
kind: Pod
metadata:
  name: opencl-cpu
spec:
  containers:
    - name: opencl-container
      image: rsankar12/opencl_cse160: latest
      command: ["sleep", "infinity"]  # Keeps the pod running for testing
```

This YAML file defines a Kubernetes **Pod** that sets up an environment for running and testing CSE160 assignments using OpenCL.

### Pod Configuration

- Creates a **Pod** named `opencl-cpu`.
- Inside the pod, it launches one **container** named `opencl-container`.
- Uses the Docker image `rsankar12/opencl_cse160:latest`, a modified version of the POCL container originally used in Professor Kastner’s CSE160 course.
- The command `["sleep", "infinity"]` ensures the container stays alive indefinitely, making it ideal for interactive testing and debugging.

### What’s Included in the Image**

- **POCL (Portable OpenCL)**: Provides OpenCL support for CPU-based execution.
- **`helper_lib/`**: Contains pre-installed helper libraries required by the assignment.
- **`Makefile`**: A customized Makefile adapted to run in the Kubernetes environment.
- **Dataset and Testing Files**: Sample datasets and test files are included for immediate use and evaluation.

To run the assignment inside the container, you need to copy the following files into the container's working directory:

- `main.c`
- `vector_add_2.cl`
- `vector_add_4.cl`

```
kubectl cp main.c opencl-cpu:/main.c
kubectl cp vector_add_2.cl opencl-cpu:/vector_add_2.cl
kubectl cp vector_add_4.cl opencl-cpu:/vector_add_4.cl
```
---
## Running CSE 160 Assignment

This Makefile automates running tests, handles errors, and formats results for GradScope grading.

### Key Features

- **Automated Testing:**  
  Runs all 10 test datasets in a loop instead of manual repetition.

- **Error Handling:**  
  Checks each test's exit status to determine pass/fail.

- **GradScope-Compatible Output:**  
  Generates a `results.json` file with detailed test results in the format GradScope expects.

- **Logging:**  
  Saves individual test outputs to `run_logs/output_i.log` for debugging.

- **Summary Reporting:**  
  Provides a human-readable grading summary with total passes and fails.

### How to Use

- Build and run all tests, generate results, and print JSON output:  
  ```
  bash
  make run
  ```

