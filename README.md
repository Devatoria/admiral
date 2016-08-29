# admiral
Admiral is a Docker Registry administration and authentication system. It is under development and is aiming to be a real production tool.

This tool is listening to Docker Registry events using notifications in order to catch `pull` and `push` events (to make some audit, for example), but also to maintain a list of available images. Next, you will be able to create teams and users, and allow only certain user to pull or push certain images.

Admiral can synchronize itself with an existing registry by getting and parsing the catalog.

## Configuration
### Configure the daemon
The configuration file of the Admiral daemon is really easy.

```
address = "127.0.0.1" # Admiral API listening address
port = 3000 # Admiral API listening port

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

### Configure the Docker Registry
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

## Synchronize the Docker Registry with Admiral

Admiral can run a set of jobs (in addition of the default daemon behavior). Admiral can synchronize itself with the Docker Registry by getting and parsing the catalog, inserting non-existing namespaces and images into the database. Just run the `synchronize` job:

```
admiral job run synchronize
```

You can have the list of available jobs by using the following command:

```
admiral job list
```
