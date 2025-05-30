# Custom image and job creation

# Building and pushing a docker image

To login

```bash
docker login -u <username>
```

Since the phones are arm64 and automaton is amd64, we need to cross build

First create a builder

```bash
docker buildx create --name myBuilder --driver docker-container --platform linux/arm64
```

Use the builder

```bash
docker buildx use myBuilder
```

(Optional) confirm that the builder is configured correctly

```bash
docker buildx inspect --builder myBuilder
```

Go to the location of the Dockerfile and build and push. “ajayrr/test-py-kubernetes” is Ajay’s public dockerhub repo

```bash
docker buildx build --platform linux/arm64 -t ajayrr/test-py-kubernetes:arm64 --push .
```

A test has already been pushed and can be used directly to create a job

# Creating a job

Go to control plane node and create a yaml file for the job. (or pull from github repo, path - container-images/python-hello-world/test_py_job_helloWorld.yaml)

```bash
apiVersion: batch/v1
kind: Job
metadata:
  name: test-py-kubernetes
spec:
  template:
    spec:
      containers:
      - name: python
        image: ajayrr/test-py-kubernetes:arm64
      restartPolicy: Never
  backoffLimit: 2
```

Create a job

```bash
kubectl apply -f test_py_job_helloWorld.yaml
```

Since the job only prints hello world it completes quickly. The job pod can be found by viewing all pods

```bash
kubectl get pods --all-namespaces
```

To view logs

```bash
kubectl logs <podname> -n <namespace>
```

If the job was successful, the logs will contain “Hello from Kubernetes Job!!”