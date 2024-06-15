## Overview

## Apps

### Frontend

Minimal Golang HTTP App that serves HTML with name day.
This data is fetched from the backend server

Expects env variable for declaring backend hostname

```env
BACKEND_URL="http://backend:3001"
```

### Backend

Again, minimal Golang app that servers HTTP REST API. It connects to Czech API for current name day info and caches it.

### Copy of K8s Install on Fedora

Deployed for a demo on a Fedora 40 ARM running K8s. Edited
[Fedora Guide](https://docs.fedoraproject.org/en-US/quick-docs/using-kubernetes/#sect-fedora40-and-newer) so that I have 2 nodes running.

#### Node-1

```bash
echo "node-1 > /etc/hostname"
sudo systemctl disable --now firewalld
sudo dnf install iptables iproute-tc
sudo cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
sudo modprobe overlay
sudo modprobe br_netfilter

sudo cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
sudo sysctl --system

lsmod | grep br_netfilter
lsmod | grep overlay
sysctl net.bridge.bridge-nf-call-iptables net.bridge.bridge-nf-call-ip6tables net.ipv4.ip_forward

sudo dnf install cri-o containernetworking-plugins
sudo dnf install kubernetes kubernetes-kubeadm kubernetes-client
sudo systemctl enable --now crio
sudo systemctl enable --now kubelet
```

Duplicate the VM as node-2

#### Node-1

```bash
sudo kubeadm init --pod-network-cidr=10.244.0.0/16

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

kubectl taint nodes --all node-role.kubernetes.io/control-plane-
kubectl apply -f https://github.com/coreos/flannel/raw/master/Documentation/kube-flannel.yml
kubeadm token create --print-join-command
```

#### Node-2

```bash
echo "node-2 > /etc/hostname"
sudo kubeadm join 192.168.67.6:6443 --token m5k7x8......
```

#### Node-1

```bash
user@node-1:~$ kubectl get nodes
NAME     STATUS   ROLES           AGE     VERSION
node-1   Ready    control-plane   3m49s   v1.29.5
node-2   Ready    <none>          17s     v1.29.5
```

### Running apps

```bash
kubectl apply -f backend.yaml
kubectl apply -f frontend.yaml
kubectl get services
```
