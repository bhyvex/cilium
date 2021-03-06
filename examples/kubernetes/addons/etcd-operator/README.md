Kubernetes etcd-operator integration (beta)
===========================================

This directory contains the necessary scripts to deploy etcd-operator with Cilium in your development cluster.

Prerequistes
------------

Ensure the dependencies are installed. The CloudFlare PKI/TLS toolkit is leveraged for generating certificates. The deployment below requires `cfssl` and `cfssljson` which you can download [here](https://pkg.cfssl.org/). Ensure the binaries are in your `PATH`. 

Alternatively, if you have Go installed, you can get the libraries with the following:

```
go get -u github.com/cloudflare/cfssl/cmd/cfssl
go get -u github.com/cloudflare/cfssl/cmd/cfssljson
```

Create etcd certificates
------------------------

The first step you need to do is

```
tls/certs/gen-cert.sh <cluster domain>
```
where `<cluster domain>` is the domain the cluster set up in kube-dns.

You can find it by checking the config map of core-dns by running
```
kubectl get ConfigMap --namespace kube-system coredns -o yaml | grep kubernetes
```

or by checking the kube-dns deployment and grepping for 'domain'
```
kubectl get Deployment --namespace kube-system kube-dns -o yaml | grep domain
```

_Note: Make sure to remove the trailing dot from the domain set up in
the options. For example, if the flag option is `cluster.local.` the cluster
domain for the `gen-cert.sh` should be `cluster.local`_

For reference, the cluster domain used in Kubernetes clusters by default is `'cluster.local'`

Deploy generated certificates
-----------------------------

The next step is to deploy the certificates generated by the previous step
in the Kubernetes cluster, this can be achieved by running:

```
tls/deploy-certs.sh
```

Deploy kube-dns and make sure it contains the required label
------------------------------------------------------------

Ensure your DNS service is running and contains the label `io.cilium.fixed-identity=kube-dns` using the following command.

# if using kube-dns
```
kubectl patch -n kube-system deployment/kube-dns --type merge -p '{"spec":{"template":{"metadata":{"labels":{"io.cilium.fixed-identity":"kube-dns"}}}}}'
```
# if using coredns
```
kubectl patch -n kube-system deployment/coredns --type merge -p '{"spec":{"template":{"metadata":{"labels":{"io.cilium.fixed-identity":"kube-dns"}}}}}'
```

Deploy Kubernetes descriptors for etcd operator as well Cilium
--------------------------------------------------------------

Just simply run `kubectl create -f` on this directory.

```
kubectl create -f ./
```

Please wait until the `etcd-operator` and `cilium-etcd-cluster.yaml` pods are in
ready state.

Wait a couple seconds and everything should be running fine. All pods, including
Cilium can suffer from restarts until the system converge to a readiness state.
This is expected for the first 3 to 5 minutes.

```
$ kubectl get pods -n kube-system
NAME                             READY     STATUS    RESTARTS   AGE
cilium-etcd-g2sr9qxdhw           1/1       Running   0          6h
cilium-etcd-ss5jlv4cbq           1/1       Running   0          7h
cilium-etcd-x28h2rkhz7           1/1       Running   0          7h
etcd-operator-69b5bfc669-fvm8c   1/1       Running   0          1d
kube-dns-7dcc557ddd-tsqjt        3/3       Running   12         1d
```
