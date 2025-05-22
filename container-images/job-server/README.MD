To compile the go binary use `GOOS=linux GOARCH=arm GOARM=7 go build -o jobserver -ldflags="-s -w" server.go`

To build the docker image for the phones use `docker buildx build --platform linux/arm64 -t arunanthivi/k8s-job-server:latest . --push` and restart the deployment from the control plane