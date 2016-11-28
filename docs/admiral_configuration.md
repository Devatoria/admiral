# Admiral configuration

## Introduction

Admiral can be configured using a `config.toml` file locating at:

* `/etc/admiral` (production use)
* in the same directory than the binary (development use)

Here is an example:

```
debug = true

address = "127.0.0.1"
port = 3000

[auth]
issuer = "pepito"
token-expiration = 5
certificate = "/certs/server.crt"
private-key = "/certs/server.key"

[admin]
username = "admin"
password = "password"

[database]
host = "localhost"
port = 5432
user = "admiral"
password = "admiral"
name = "admiral"

[registry]
address = "http://localhost"
port = 5000
```

## General configuration

#### debug

It enables the Gin framework debug mode (printing routes and some things).

#### address & port

Address (interface) and port admiral should listen to.

## Authentication configuration

#### issuer

This parameter is shared between the docker registry and admiral and must be the same, otherwise, authentication will not work.

#### token-expiration

Bearer token expiration in minutes. After this time, the authentication will need another token.

#### certificate & private-key

Certificate and private key used to sign bearer token (and used to secure the registry).

## Registry configuration

The registry configuration is used to allow admiral to contact the registry, for example, to delete images.

## Admin configuration
Set an admin user which be able to do action on all images
