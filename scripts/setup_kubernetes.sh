#!/bin/sh

sudo apk update
sudo apk del iptables
sudo apk add --no-interactive ethtool
sudo apk add --no-interactive iproute2
sudo apk add --no-interactive socat

sudo service init_starling start
sudo swapoff -a

iptables -D KUBE-FIREWALL 1
sudo kubeadm reset --force

sudo hwclock --systohc
sudo umount -f /var
sudo mount -o bind /var/lib/kubelet /var/lib/kubelet
sudo mount --make-shared /var/lib/kubelet
rm $HOME/.kube/config
