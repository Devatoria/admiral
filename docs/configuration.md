# Configuration

## Admiral

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
