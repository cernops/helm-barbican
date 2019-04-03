# Helm Barbican Secrets Plugin

## Installation

```bash
$ helm plugin install https://gitlab.cern.ch/helm/plugins/barbican
```

## Secrets

Barbican is used for secret storage, first step is to get an openstack token.

```bash
$ export OS_TOKEN=$(openstack token issue -c id -f value)
$ export OS_AUTH_URL=https://keystone.cern.ch/v3
$ unset OS_IDENTITY_PROVIDER OS_AUTH_TYPE OS_MUTUAL_AUTH OS_PROTOCOL
```

Typical usage will be `edit` to change your secrets and `view` to display them.

```
$ helm secrets edit secrets.yaml
param1:
  subparam2: value2
  subparam3: value3

$ helm secrets view secrets.yaml
param1:
  subparam2: value2
  subparam3: value3
```

The plugin provides wrapper commands for `install`, `upgrade` and `lint` and
handles the secrets transparently - they are decrypted into shared memory and
passed to helm being deleted right after.

```
$ helm secrets install stable/nginx --name nginx --namespace nginx --values secrets.yaml

$ helm secrets upgrade nginx stable/nginx --values secrets.yaml

$ helm secrets lint stable/nginx --values secrets.yaml
```

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

## Development

The plugin is a go binary. It requires go>=1.11 (no vendor or GOPATH needed).
```
go build .
```

To reinstall the plugin locally for testing:
```
helm plugin remove secrets; helm plugin install `pwd`
```
