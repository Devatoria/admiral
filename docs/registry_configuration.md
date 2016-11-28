# Registry configuration

## Introduction

The Docker Registry needs to be configured in order to complete two tasks:

* send events for audit
* authenticate

All of these configurations must be in the `/etc/docker/registry/config.yml` file.

## Events

This will ask to the registry to send events to Admiral. It will store these events in a database in order to allow to make some audit on this. It will also get push events in order to maintain images database up to date.

```
notifications:
  endpoints:
  - name: admiral
    disabled: false
    url: http://<admiral_host>:<admiral_port>/events
    timeout: 500ms
    threshold: 5
    backoff: 5s
```

## Authentication

This will enable authentication. Now, the registry will ask to admiral for token based authentication (for example, if an user want to pull an image).

The `issuer` and the `rootcertbundle` must be shared between admiral and registry, otherwise, the registry will not be able to verify tokens.

```
auth:
  token:
    realm: http://<admiral_host>:<admiral_port>/v1/token
    service: registry
    issuer: <admiral_issuer>
    rootcertbundle: <admiral_private_key>
```
