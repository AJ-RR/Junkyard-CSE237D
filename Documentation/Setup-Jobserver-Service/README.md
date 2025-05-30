# GreenGrader JobServer Service

Files & media: server.go

I set up a pod to help simplify the process of creating jobs for grading assignments.

The code is attached to this document. If modifications are to be made, be sure to precompile the code with the command

`GOOS=linux GOARCH=arm GOARM=7 go build -o jobserver server.go` 

to ensure that the executable properly runs on the phones. 

## Compiling the Server

The go code for the job-server is attached to this note and included in the project repo. To adapt the code to different assignments, it only requires a change to two lines (:135)

```go
Image:"<GRADING IMAGE>"
ImagePullPolicy: corev1.PullAlways,
Command:[]string{"sh", "-c", "unzip /scripts/archive.zip -d $HOME >/dev/null 2>&1 && <TEST COMMAND>"},
```

The server will pull the defined image (which should include all dependencies and harness files), unzip the student’s submission into `$HOME` directory, and run the specified command from `$HOME` to generate the testing results. 

**NOTE**: The test command must **ONLY** write the results.JSON information to `stdout`, all other output must be suppressed

**NOTE:** The current implementation of the job server uses a ConfigMap to pass in the submission into the job. Because of the limitations of Kubernetes, this means that the zipped submission files **CANNOT** be ≥ 1MiB. Anything not unique to a student’s submission should be baked into the image.

To build the binary for the phones (requires golang)

```bash
GOOS=linux GOARCH=arm GOARM=7 go build -o jobserver -ldflags="-s -w" server.go
```

To initialize the binary on the phones, we use a simple Dockerfile to define our container

```docker
FROM alpine:3.21

WORKDIR /app
COPY jobserver .

RUN chmod +x ./jobserver

EXPOSE 5000
ENTRYPOINT [ "./jobserver" ]
```

To update the docker image with the new recompiled binary, use

```bash
docker buildx build --platform linux/arm64 -t <YOUR DOCKER REPO>:latest . --push
```

## Creating the Server

We start by creating a new YAML file for the server `jobserver.yaml`

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: job-server-sa
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: job-server-role
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["configmaps"]
    verbs: ["create"]
  - apiGroups: ["batch"]
    resources: ["jobs"]
    verbs: ["create", "get"]
  - apiGroups: [""]
		resources: ["pods"]
		verbs: ["list", "get"]
  - apiGroups: [""]
		resources: ["pods/log"]
		verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: job-server-binding
  namespace: default
subjects:
  - kind: ServiceAccount
    name: job-server-sa
    namespace: default
roleRef:
  kind: Role
  name: job-server-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: job-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: job-server
  template:
    metadata:
      labels:
        app: job-server
    spec:
      serviceAccountName: job-server-sa
      containers:
        - name: job-server
          image: <YOUR DOCKER REPO>:latest
          ports:
            - containerPort: 5000
---
apiVersion: v1
kind: Service
metadata:
  name: job-server-service
spec:
  type: NodePort
  selector:
    app: job-server
  ports:
    - port: 5000
      targetPort: 5000
      protocol: TCP
      nodePort: 30080
```

This is the full YAML file for the job server. It is separated into 3 main parts

- The Deployment pulls the image (the precompiled go binary) from my dockerhub repo and exposes the port 5000 on the container
- The Service uses NodePort (same as the [Setting up Web Server](https://www.notion.so/Setting-up-Web-Server-1d686fd69b8a80ff9a26e6f9be5c8d3b?pvs=21)) to connect the container’s 5000 port to the 30080 port (this is the port number we actually use to send requests)
- This container requires seperate permissions to create ConfigMaps (for the scripts) and to create jobs, so we have to create a new ServiceAccount for the deployment. we call `job-server-sa` . Alongside creating this account, we have to attach the relevant roles and permissions to this account

With this file on the cluster, we should be able to get the server up and running with

```bash
kubectl apply -f jobserver.yaml
```

You can make sure the pod is running by checking the pod with `kubectl get pods` and `kubectl get svc`. 

The former should have some pod running with the name `job-server-xxxxx-xxxxx` . 

## Testing the Server from outside the Cluster

### From the Cluster:

To access the server, we’ll need the IP of the nodes in our cluster. We can retrieve these with the command `kubectl get nodes -o wide`

My cluster gives me a result like this

```bash
NAME             STATUS   ROLES           AGE     VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE            KERNEL-VERSION               CONTAINER-RUNTIME
google-felix-5   Ready    control-plane   4d22h   v1.29.6   10.42.0.6     <none>        postmarketOS edge   5.10.198-android13-4-dirty   containerd://1.6.36
google-felix-6   Ready    <none>          4d22h   v1.29.6   10.42.0.7     <none>        postmarketOS edge   5.10.198-android13-4-dirty   containerd://2.0.5
```

With my NodePort service already running on port 30080, my url to access the server becomes 

**http://10.42.0.[6/7]:30080**

### From Automaton:

Test that you can correctly access the server by just pinging the host route

```bash
curl http://10.42.0.7:30080
```

If successful, this should return the response

```bash
{"message":"Server is running"}
```

Now that we can successfully access the server from outside the cluster, we can attempt to send a job.

## Creating an Assignment

Currently the job server only creates a new job based on a python image and runs the included python file. We have to first create a python file to send as our assignment.

```bash
echo 'print("Hello from the kubernetes cluster!")' >> script.py
```

With our assignment file made, we can create our request to the cluster with curl

```bash
curl -X POST http://10.42.0.7:30080/submit -F "name=<STUDENTNAME>" -F "image=<ASSIGNMENTNAME>" -F "script=@<path to script.py>"
```

(NOTE: make sure that the name and image parameters are all lowercase or else the script will fail)

(If `script.py` exists in our current directory, it’s path is just `@script.py` )

If successful, this should return the response

```bash
{"status":"Job created"}
```

## Viewing the Assignment Result

### From the Cluster:

With the successful job creation, we can see its status with `kubectl get jobs` 

This would return some output like 

```bash
NAME                      COMPLETIONS   DURATION   AGE
testuser-testassignment   1/1           25s        36m
```

(TODO: the `server.go` code likely needs to be recompiled to include some unique identifier appended to the end of the job name. Current implementation likely fails upon re-submission)

I can check the stdout output of this job with `kubectl logs job/testuser-testassignment` 

```bash
Hello from the kubernetes cluster!
```