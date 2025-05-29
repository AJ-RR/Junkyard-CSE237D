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
    namespace   = "test-assignments-git"
    image       = "ajayrr/opencl-kube-git:arm64"

    config.load_kube_config(config_file=kubeconfig)
    batch_v1 = client.BatchV1Api()

    total_sec = 0.0

    for i in range(n):
        job_name = f"grader-job-{i+1}-{datetime.now().strftime('%H%M%S')}"
        job_spec = create_job_template(job_name, namespace, image)

        # 1) submit the Job
        start = time.time()
        batch_v1.create_namespaced_job(namespace=namespace, body=job_spec)
        print(f"[✓] Submitted {job_name}")

        # 2) poll until the Job completes (Succeeded or Failed)
        while True:
            job = batch_v1.read_namespaced_job(name=job_name, namespace=namespace)
            if job.status.succeeded and job.status.succeeded > 0:
                break
            if job.status.failed and job.status.failed > 0:
                print(f"[✗] {job_name} failed")
                break
            time.sleep(2)

        # 3) record latency for this Job
        delta = time.time() - start
        total_sec += delta
        print(f"    ↳ {job_name} latency: {delta:.2f}s")

        time.sleep(delay)          # optional spacing before next submit

    print(f"\nTOTAL latency for {n} jobs: {total_sec:.2f}s ({total_sec/60:.2f} min)")

