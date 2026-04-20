# How to bootstrap a new node on AWS

The following guide explains how to bootstrap a new node running k3s in AWS.
The purpose is so that the application can be deployed to this node.

## Assign Elastic IP

Assign the elastic IP to the instance first.

## Preparing AWS SSH Key

After downloading the SSH Key, it goes in `./local` directory.

```sh
cd ./local
chmod 400 "<key>.pem"
ssh -i "<key>.pem" ubuntu@<ip>
```

## Installing k3s

1. Follow these steps [to disable the firewall](https://docs.k3s.io/installation/requirements?os=debian)

2. Install via this command:

```sh
curl -sfL https://get.k3s.io | sh -
```

## Configuring k3s

After assigning the elastic IP to the instance, add an additional SAN:

```sh
sudo vim /etc/rancher/k3s/config.yaml
```

Add the following to the file:

```yaml
tls-san:
  - "<ip>"
```

Reboot the service:

```sh
sudo systemctl restart k3s.service
```

## Connecting to the k3s cluster

Get the k3s.yaml file:

```sh
ssh -i "<key>.pem" ubuntu@<ip> sudo cat /etc/rancher/k3s/k3s.yaml > ./server.yaml
```

Replace the ip in the file with the ip of the server:

```sh
sed -i -e 's/server: https:\/\/127.0.0.1:6443/server: https:\/\/<ip>:6443/' server.yaml
```

Test the connection:

```sh
KUBECONFIG=server.yaml kubectl get pods -A

# output:
kube-system   coredns-76c974cb66-w42b5                  1/1     Running     0          16m
kube-system   helm-install-traefik-49zb4                0/1     Completed   2          16m
kube-system   helm-install-traefik-crd-p8f87            0/1     Completed   0          16m
kube-system   local-path-provisioner-8686667995-5mfkn   1/1     Running     0          16m
kube-system   metrics-server-c8774f4f4-wj6vk            1/1     Running     0          16m
kube-system   svclb-traefik-c6772f2a-6xmpl              2/2     Running     0          16m
kube-system   traefik-c5c8bf4ff-t7qll                   1/1     Running     0          16m
```