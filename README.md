# `operator-tools`

CLI for me

```
go get github.com/tlwr/operator-tools
alias op='operator-tools'
```

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

### HTTP

```
op http
```

Interact with
[CredHub](https://docs.cloudfoundry.org/credhub/).

Utilities for doing things via HTTP

### Profile a HTTP request/response

Produces a timeline of an HTTP request: DNS, TCP, TLS, request, response.

```
op http profile -u https://healthcheck.cloudapps.digital/
|=============================================================================================| total duration 309ms
|~~~~~~~~~~~~                                                                                 | dns from 0ms until 43ms duration 43ms
|            ~~~~                                                                             | connect from 43ms until 59ms duration 16ms
|                 ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~                                        | tls from 59ms until 181ms duration 122ms
|                                                      x                                      | request-headers-done at 181ms
|                                                      x                                      | request-done at 181ms
|                                                                    ~~~~~~~~~~~~~~~~~~~~~~~~ | reading-response from 227ms until 309ms duration 82ms
```

### YAML

```
op yaml
```

Interact with
[CredHub](https://docs.cloudfoundry.org/credhub/).

Do things with [YAML](https://yaml.org/).

### Find things in YAML

Use [BOSH interpolate's path syntax](https://bosh.io/docs/cli-int/) to traverse
YAML.

```
# Access by key
op yaml find -p /users < f.yml

# Access by key and value
op yaml find -p /users/name=tlwr < f.yml

# First element of a list
op yaml find -p /0 < f.yml

# Access by key and value, and list element
op yaml find -p /users/name=tlwr/roles/0 < f.yml
```
