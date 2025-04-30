from kubernetes import client, config, utils
import yaml
import os

def create_kubernetes_job(job_name, job_yaml_path, kubeconfig_path, namespace="default"):
    """
    Creates a Kubernetes job in a specified namespace using a specified kubeconfig file.

    Args:
        job_name (str): The name of the Kubernetes job to create.
        job_yaml_path (str): The path to the YAML file defining the job.
        kubeconfig_path (str): The path to the kubeconfig file.
        namespace (str, optional): The namespace in which to create the job. Defaults to "default".
    """
    try:
        # 1. Load Kubernetes configuration from the specified file
        config.load_kube_config(config_file=kubeconfig_path)

        # 2. Create a Kubernetes API client for core (batch) operations
        batch_v1_api = client.BatchV1Api()

        # 3. Load the job definition from the YAML file
        with open(job_yaml_path, 'r') as f:
            job_definition = yaml.safe_load(f)
            #Important: Set the metadata name here, before creating.
            job_definition['metadata']['name'] = job_name

        # 4. Create the Kubernetes job
        api_response = batch_v1_api.create_namespaced_job(
            body=job_definition,
            namespace=namespace
        )

        print(f"Job '{job_name}' created successfully in namespace '{namespace}'.")
        print(f"Response: {api_response}")

    except Exception as e:
        print(f"Error creating job '{job_name}': {e}")
        raise  # Re-raise the exception for further handling if needed


def main():
    """
    Main function to define job parameters and call create_kubernetes_job.
    """
    job_name = "test-py-kubernetes"
    job_yaml_path = "test_job.yaml"
    kubeconfig_path = "fold_1_config.conf"  #  Path to your kubeconfig file
    namespace = "default"
    try:
        create_kubernetes_job(job_name, job_yaml_path, kubeconfig_path, namespace)
    except Exception as e:
        print(f"Failed to create kubernetes job: {e}")


if __name__ == "__main__":
    main()