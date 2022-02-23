# wellknown-matrix

A Docker container created to easily serve Matrix's Federation and Client .well-known file, without
fiddling too much with configuration at your reverse proxy.

## Usage
> WARNING: this container is meant to be used behind a reverse proxy.

> DO NOT DIRECTLY EXPOSE TO THE INTERNET.

If you're using docker-compose, you can add something like this to your services:
```
    container_name: wellknown-matrix
    image: rikhwanto/wellknown-matrix:0.0.1
    restart: unless-stopped
    environment:
      - FEDERATION_SERVER=example.com:443
      - CLIENT_HOMESERVER=https://example.com
      - CLIENT_IDENTITYSERVER=https://identity.com
```
It uses port 8080 by default and you can forward requests to `.well-known/matrix/` from your reverse proxy to this container.

### ENVIRONMENT VARIABLES
`FEDERATION_SERVER` is the address where your Matrix federation server is located and will be served at `/.well-known/matrix/server`. You can read more about this server-server API [here](https://matrix.org/docs/spec/server_server/r0.1.0#server-discovery)

`CLIENT_HOMESERVER` is the URL of homeserver for the Matrix client to connect to and `CLIENT_IDENTITYSERVER` (not required, you can skip this) is the URL of optional identity server for the Matrix client to connect to. They will be served at `/.well-known/matrix/client`. You can read more about this client-server API [here](https://matrix.org/docs/spec/client_server/r0.4.0#server-discovery)

### REVERSE PROXY SETUP
This is an example of how to setup the reverse proxy using Nginx. You can add something like this to your Nginx configuration:
```
    # other configuration

    location /.well-known/matrix {
        proxy_pass http://wellknown-matrix:8080/.well-known/matrix
    }

    # other configuration
```
