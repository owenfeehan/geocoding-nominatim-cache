# geocoding-nominatim-cache

## Author

[Owen Feehan](http://www.owenfeehan.com)

## Description

A Go ([GIN](https://gin-gonic.com/)) service that provides an RESTful end-point for geocoding placenames via [Nominatim](https://nominatim.org/) (geocoding service for [OpenStreetMap](https://www.openstreetmap.org/)).

Requests to Nomatim are throttled and cached, as required by the [Nominatim Usage / Geocoding Policy](https://operations.osmfoundation.org/policies/nominatim/). The cache-storage uses either a local [BadgerDB](https://github.com/hypermodeinc/badger) persisent file (default) or a [Redis](https://redis.io/) backend.

The service allows a database of geolocated data to be built up over time within a particular environment (e.g. corporate or personal). If necessary, the Redis backend can be configured to expire data.

The RESTful end-point is compliant with Swagger/OpenAPI. See `http://localhost:8080/swagger/index.html` (or whatever address the service becomes bound to) and `http://localhost:8080/swagger/doc.json`. The [OpenAPI generator](https://github.com/OpenAPITools/openapi-generator) can quickly create an automated client across many languages and frameworks.

It requires Go v1.21 at a minimum.

## License

MIT, see `LICENSE.txt`.

## Usage

Start the service with either:

> go run main.go --debug

or by building a binary with:

> go build -o geocoding-nominatim-cache main.go --debug

The `--debug` arugment can be dropped for production use.

### CLI Arguments

| Argument            | Type     | Default                 | Description                                                                                                                                             |
|---------------------|----------|-------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--address`         | string   | `localhost:8080`        | The address to bind the server to. See [Gin.Run](https://pkg.go.dev/github.com/gin-gonic/gin#Engine.Run).                                               |
| `--redis`           | string   | *use BadgerDB instead*  | Binds to a redis server at the given address (e.g., `localhost:6379`). If not set or empty, uses BadgerDB as the default store.                           |
| `--debug`           | bool     | `false`                 | Enable debug logging and debug mode on the web server.                                                                                                  |
| `--inMemory`        | bool     | `false`                 | Uses in-memory (non-persistent) location storage, ignoring Redis or BadgerDB. This takes precedence over the redis flag, if set.                        |
| `--throttle`        | int      | `2000`                  | The minimum number of milli-seconds between requests to the Nominatim API (at least 1000 milliseconds are required by [policy]((https://operations.osmfoundation.org/policies/nominatim/)). |
| `--trusted-proxies` | string   | *trust no IPs*          | Comma-separated list of trusted proxy IPs or CIDRs. See Gin's [SetTrustedProxies](https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies) |
