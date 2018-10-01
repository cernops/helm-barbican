# Helm Deployment and Secrets Manager

## Installation

```bash
$ helm plugin install https://gitlab.cern.ch/helm/plugins/cern
```

## Secrets

Barbican is used for secret storage, first step it to get an openstack token.

```bash
$ export OS_TOKEN=$(openstack token issue -c id -f value)
$ export OS_AUTH_URL=https://keystone.cern.ch/v3
$ unset OS_IDENTITY_PROVIDER OS_AUTH_TYPE OS_MUTUAL_AUTH OS_PROTOCOL
```

Typical usage will be 'edit' to change your secrets and 'view' to display them.

```
$ helm cern edit
param1:
  subparam2: value2
  subparam3: value3

$ helm cern view
param1:
  subparam2: value2
  subparam3: value3
```

Commands `enc` and `dec` offer lower level functionality to encode and decode
the secrets.yaml file, but you should not usually need them.
