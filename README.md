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

## Development

The plugin is a go binary. It requires go>=1.11 (no vendor or GOPATH needed).
```
go build .
```

To reinstall the plugin locally for testing:
```
helm plugin remove secrets; helm plugin install `pwd`
```
