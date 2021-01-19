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

#### Visualize certificates

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

Utilities for doing things via HTTP

#### Profile a HTTP request/response

Produces a timeline of an HTTP request: DNS, TCP, TLS, request, response.

```
op http profile -u https://healthcheck.cloudapps.digital/
TLS:
Version: TLSv1.2
Cipher-Suite: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
Server-Name: healthcheck.cloudapps.digital
Negotiated-Protocol: h2

Status:
HTTP/2.0 200 OK

Headers:
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Vcap-Request-Id: 0b8b6360-a978-4dc2-645a-404c5951c194
Date: Tue, 19 Jan 2021 15:17:42 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 222223
Accept-Ranges: bytes
Cache-Control: max-age=0,no-store,no-cache
Last-Modified: Tue, 19 Jan 2021 12:40:10 GMT

Trace:
|=======================================================| total duration was 607ms
|~~~~~~~~~~                                             | dns from 0ms until 139ms duration 139ms
|          ~~~~                                         | connect from 140ms until 174ms duration 34ms
|              ~~~~~~~~~~~~~~~~~~~~~~~~~                | tls from 174ms until 432ms duration 258ms
|                                        x              | request-headers-done at 432ms
|                                        x              | request-done at 432ms
|                                             ~~~~~~~~~ | reading-response from 479ms until 607ms duration 128ms
```

### x509

```
op x509
```

Do things with x509 (eg certificates).

#### Find certificates

```
# Finds expiring certificates
op x509 fc

# Finds certificates expiring within 180 days
op x509 fc --expiry-days 180

# Finds certificates excluding fixtures and testdata directories
op x509 fc --exclude fixtures,testdata
```



### YAML

```
op yaml
```

Do things with [YAML](https://yaml.org/).

#### Find things in YAML

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
