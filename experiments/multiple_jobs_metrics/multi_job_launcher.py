from kubernetes import client, config
import yaml
import time
from datetime import datetime

def create_job_template(job_name, namespace, image):
    command = (
        "date; "
        "git clone https://github.com/AJ-RR/Dummy-CSE160-Assignment.git && "
        "cd Dummy-CSE160-Assignment/CSE160Assignment2 && "
        "time ./evaluate_assignment.sh > output.log 2>&1; "
        "date; "
        "cat output.log"
    )

    return {
        "apiVersion": "batch/v1",
        "kind": "Job",
        "metadata": {
            "name": job_name,
            "namespace": namespace
        },
        "spec": {
            "template": {
                "spec": {
                    "containers": [{
                        "name": f"{job_name}-container",
                        "image": image,
                        "imagePullPolicy": "Always",
                        "command": ["/bin/bash", "-c", command]
                    }],
                    "restartPolicy": "Never"
                }
            },
            "backoffLimit": 0
        }
    }

def submit_job(kubeconfig_path, job_spec, namespace):
    config.load_kube_config(config_file=kubeconfig_path)
    batch_v1 = client.BatchV1Api()

    try:
        batch_v1.create_namespaced_job(
            namespace=namespace,
            body=job_spec
        )
        print(f"[✓] Submitted job: {job_spec['metadata']['name']}")
    except Exception as e:
        print(f"[✗] Failed to submit job {job_spec['metadata']['name']}: {e}")

def launch_multiple_jobs(n=5, delay=2):
    kubeconfig = "fold_1_config.conf"
    namespace = "test-assignments-git"
    image = "ajayrr/opencl-kube-git:arm64"

    for i in range(n):
        timestamp = datetime.now().strftime("%H%M%S")
        job_name = f"grader-job-{i+1}-{timestamp}"
        job_spec = create_job_template(job_name, namespace, image)
        submit_job(kubeconfig, job_spec, namespace)
        time.sleep(delay)  # optional delay to avoid burst creation

if __name__ == "__main__":
    launch_multiple_jobs(n=5, delay=2)
