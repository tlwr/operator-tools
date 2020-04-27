# `operator-tools`

CLI for me

## Commands

### CredHub

```
op credhub
```

Interact with
[CredHub](https://docs.cloudfoundry.org/credhub/).

### Visualize certificates

Produces a particularly illegible DAG of CredHub certificates and certificate
authorities.

```
credhub curl -p /api/v1/cerficiates  |\
opt credhub visualize-certificates |\
dot -Tpng -o /tmp/certs.png ;
open /tmp/certs.png
```
