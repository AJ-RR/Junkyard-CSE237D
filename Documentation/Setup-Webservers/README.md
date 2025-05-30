# Setting up Web Server

From **automaton**, ssh into the control plane of the cluster (e.g. 10.42.0.2)

Create the deployment file for NGINX:

**NGINX Container file**

```
yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
```

With the deployment file created, create the container with:

```
kubectl create -f [PATH TO nginx-deployment.yaml]
```

Verify that the container is up correctly with

```
kubectl get deployments
```

Create the service file for NGINX

NGINX Service file

```
yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  type: NodePort
  selector:
    app: nginx
  ports:
    - port: 80
      targetPort: 80
      nodePort: 30080
```

With the service file created, start the service with:

```
kubectl apply -f [PATH TO nginx-service.yaml]
```

Confirm that the NGINX service is running with

```
kubectl get svc
```

In the output, there should be a row with the name

```
nginx-service
```

Confirm that the web server is running

```
curl http://[WORKER IP]:[NGINX-SERVICE NODEPORT]
```

This should return the basic HTML template for NGINX

**Forwarding internal port to local machine**

```
ssh  -p44422 -L 8080:[WORKER IP]:[NGINX-SERVICE NODEPORT] [AUTOMATON USERNAME]@smartcycling.sysnet.ucsd.edu
```

After this, the page should be accessible from your browser at

```
http://localhost:8080
```