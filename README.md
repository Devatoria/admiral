# Admiral
Admiral is a Docker Registry administration and authentication system. It is under development and is aiming to be a real production tool.

It works on 2 ways:

* The first is that Docker Registry is sending events to Admiral, so it can store them in database in order to make an audit of the registry.
* The second is that Docker Registry calls Admiral to authenticate actions in order to restrict pulls and pushes to autorized user only.

Admiral can synchronize itself with an existing registry using `synchronize` job described below.

Actually, when you create a new user, the associated namespace is created. This namespace is private and personal: only the owner can push and pull on/from it.

## Features

* Existing registry synchronization
* Auto-update images and tags lists by listening to registry events
* Authenticated calls
* Private and personal namespaces
* Image deletion
* Public images

## Roadmap

Features:

* Team management
* Quota management

Side projects:

* API documentation
* CLI
* Web UI

## Configuration
### Configure the daemon
The configuration file of the Admiral daemon is really easy.

```
debug = true # Enable debug mode

address = "127.0.0.1" # Admiral API listening address
port = 3000 # Admiral API listening port

[auth]
issuer = "pepito" # Registry auth issuer, must be the same as the registry configuration
token-expiration = 5 # Token expiration time in minutes
certificate = "/certs/server.crt" # Certiciate path
private-key = "/certs/server.key" # Certificate private key path

[database]
host = "localhost" # Database host
port = 5432 # Database port
user = "admiral" # Database user
password = "admiral" # Database password
name = "admiral" # Database name

[registry]
address = "http://localhost" # Docker Registry address
port = 5000 # Docker Registry port
```

### Configure the Docker Registry authentication
In your `/etc/docker/registry/config.yml`, please add the following notification endpoint:

```
auth:
  token:
    realm: http://localhost:3000/v1/token
    service: registry
    issuer: pepito
    rootcertbundle: /certs/server.crt
```

### Configure the Docker Registry events
In your `/etc/docker/registry/config.yml`, please add the following notification endpoint:

```
notifications:
  endpoints:
  - name: admiral
    disabled: false
    url: http://admiral_host:3000/events
    timeout: 500ms
    threshold: 5
    backoff: 5s
```

For more information about notifications, please check the official Docker documentation: https://docs.docker.com/registry/configuration/#/notifications

# Admiral jobs

Admiral can run several jobs (more than just launch the API). To list them, just run the `admiral job list`. Here are some details.

## Synchronize the Docker Registry with Admiral

Admiral can run a set of jobs (in addition of the default daemon behavior). Admiral can synchronize itself with the Docker Registry by getting and parsing the catalog, inserting non-existing namespaces and images into the database. Just run the `synchronize` job:

```
admiral job run synchronize
```

# Run tests

This project contains several tests for different packages. If you want to contribute, and in order to be sure that your changes do not impact the project behavior, you can run tests using this command at the root of the project:

```
go test ./...
```
