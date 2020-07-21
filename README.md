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
HTTP/2.0 200 OK

Headers:
Accept-Ranges: bytes
Cache-Control: max-age=0,no-store,no-cache
Last-Modified: Fri, 12 Jun 2020 11:27:56 GMT
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Vcap-Request-Id: 71c47bb4-33e4-417b-4388-397596a77894
Date: Sun, 14 Jun 2020 22:50:25 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 222223

Trace:
|=============================================================================================| total duration 266ms
|~~~~~~~~                                                                                     | dns from 0ms until 24ms duration 24ms
|        ~~~~~                                                                                | connect from 24ms until 40ms duration 16ms
|             ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~                                     | tls from 40ms until 164ms duration 124ms
|                                                         x                                   | request-headers-done at 164ms
|                                                         x                                   | request-done at 164ms
|                                                                       ~~~~~~~~~~~~~~~~~~~~~ | reading-response from 205ms until 266ms duration 61ms
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
