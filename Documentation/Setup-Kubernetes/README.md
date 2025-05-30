# Setting up Kubernetes

# Packages

```bash
sudo apk del iptables
sudo apk add ethtool
sudo apk add iproute2
sudo apk add socat
```

# Remove custom settings in the rootfs

In `/etc/containerd/config.toml` , set/replace

```bash
sandbox_image = "registry.k8s.io/pause:3.9"
```

Then run

```bash
sudo rc-service containerd restart
```

# Setup (some might need sudo)

```bash
sudo service init_starling start
sudo swapoff -a

sudo iptables -D KUBE-FIREWALL 1
sudo kubeadm reset

sudo hwclock --systohc
sudo umount -f /var
sudo mount -o bind /var/lib/kubelet /var/lib/kubelet
sudo mount --make-shared /var/lib/kubelet
```

Setup the control plane. The CIDR network needs to be different per kubernetes cluster, so double check which CIDR blocks we have access to and which one to set!

```bash
sudo kubeadm init --pod-network-cidr=<CIDR network>
```

IMPORTANT: Follow the instructions on screen after the above command. It will provide a few commands that must be run to setup configs.

# Installing CNI

Currently, flannel works

```bash
wget https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
```

Open the file (`kube-flannel.yml`) and update the `net-conf.json` attribute with the CIDR provided during init. Then

```bash
kubectl apply -f kube-flannel.yml
```

To list all pods

```bash
kubectl get pods -A
```

To describe a pod

```bash
kubectl describe pod <pod name> -n <namespace>
```

To delete a pod (or restart)

```bash
kubectl delete pods <pod1> <pod2> -n <namespace>
```

To view all nodes

```bash
kubectl get node
```

To create a new join token with the command below. Use the resulting command that it gives (`kubeadm join ...` ) and run that on another phone to get it to join the cluster you setup as a node. Make sure to run setup instructions on the node before running `sudo kubeadm join ...`

```bash
sudo kubeadm token create --print-join-command
```

# Installing Helm

Helm is a package manager for Kubernetes and is required for Ray to operate through Kubernetes. To install helm, follow the steps below:

1. The phones don’t have `curl`, so using `wget` , install the arm64 tar.gz file for helm

```bash
wget https://get.helm.sh/helm-v3.17.3-linux-arm64.tar.gz
```

1. Unzip it

```bash
tar -zxvf helm-v3.17.3-linux-arm64.tar.gz 
```

1. Move the file to your `bin` folder

```bash
sudo mv linux-arm64/helm /usr/bin/helm 
```

1. Verify installation

```bash
helm version
```

# Troubleshooting

# Build missing binaries

Clone this repo and build the plugins, then copy them to the required path

[containernetworking/plugins: Some reference and example networking plugins, maintained by the CNI team.](https://github.com/containernetworking/plugins)

Delete calico related files in /etc/cni/net.d

## ImagePull Error: failed to pull image [registry.k8s.io/](http://registry.k8s.io/)…

Make sure to follow the steps in the section “Remove custom settings in the rootfs”. Rerun everything in the setup and run the command to setup the control plane again.