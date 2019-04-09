# Helm Barbican Secrets Plugin

## Installation

```bash
$ helm plugin install https://github.com/cernops/helm-barbican
```

## Secrets

Barbican is used for secret storage, first step is to get an openstack token.

```bash
export OS_TOKEN=$(openstack token issue -c id -f value)
```

Typical usage will be `edit` to change your secrets and `view` to display them.

```
helm --name mariadb secrets edit secrets.yaml
param1:
  subparam2: value2
  subparam3: value3

helm --name mariadb secrets view secrets.yaml
param1:
  subparam2: value2
  subparam3: value3
```

The plugin provides wrapper commands for helm `install`, `upgrade` and `lint`
and handles the secrets transparently - they are decrypted into shared memory
and passed to helm being deleted right after.

```
helm secrets install stable/mariadb --name mariadb --namespace mariadb --values secrets.yaml

helm secrets upgrade mariadb stable/mariadb --values secrets.yaml

helm secrets lint stable/mariadb --values secrets.yaml
```

The key used for encryption is chosen by the `--name` parameter
passed above. As an alternative if no param is passed, the cwd is used (but we
recommend relying on the helm release name).

Commands `enc` and `dec` offer lower level functionality to encode and decode
the secrets.yaml file, but you should not usually need them.

```
Available Commands:
  dec         decrypt secrets with barbican key
  edit        edit secrets
  enc         encrypt secrets with barbican key
  help        Help about any command
  install     wrapper for helm install, decrypting secrets
  lint        wrapper for helm lint, decrypting secrets
  upgrade     wrapper for helm upgrade, decrypting secrets
  view        decrypt and display secrets
```
## Kubectl plugin

You can use the secrets plugin with kubectl if you install it in your PATH as kubectl-secrets.

```
mkdir -p ~/bin

curl -o ~/bin/kubectl-secrets -L https://gitlab.cern.ch/helm/plugins/barbican/raw/master/barbican
chmod +x ~/bin/kubectl-secrets
export PATH=$PATH:$HOME/bin
```

Example usage:
```
cat secret.yaml 
/0aIzAn3AAjweH3P5Zt6nY2KPdWgVTjDrFEanpU7DjGn1/d6f3mY1JYiCTUKmSsT3xYGHj0x1XS89SWflAWXF7tmNFhgd4CMBdZNmxgc/g1tDAfheM8V1EXyJsq8aPfRHavzMFBd79C2yVMfr5wcq8PAN+knCFju5sv+QIegiZhrF6Q875X76AipQtQ=

kubectl get secret 
NAME                  TYPE                                  DATA   AGE
default-token-hjl62   kubernetes.io/service-account-token   3      6d11h

export OS_TOKEN=$(openstack token issue -c id -f value)

kubectl secrets view secret.yaml 
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm

kubectl secrets apply -f secret.yaml 
secret/mysecret created

kubectl get secret
NAME                  TYPE                                  DATA   AGE
default-token-hjl62   kubernetes.io/service-account-token   3      6d11h
mysecret              Opaque                                2      7s

kubectl get secret mysecret -o yaml
apiVersion: v1
data:
  password: MWYyZDFlMmU2N2Rm
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"password":"MWYyZDFlMmU2N2Rm","username":"YWRtaW4="},"kind":"Secret","metadata":{"annotations":{},"name":"mysecret","namespace":"default"},"type":"Opaque"}
  creationTimestamp: "2019-04-09T19:44:07Z"
  name: mysecret
  namespace: default
  resourceVersion: "1137694"
  selfLink: /api/v1/namespaces/default/secrets/mysecret
  uid: d68e4e5d-5aff-11e9-b1ff-fa163ef4a73a
type: Opaque
```

## Development

The plugin is a go binary. It requires go>=1.11 (no vendor or GOPATH needed).
```
go build .
```

To reinstall the plugin locally for testing:
```
helm plugin remove secrets; helm plugin install `pwd`
```
