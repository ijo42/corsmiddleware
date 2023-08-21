This repository includes an example plugin, `demo`, for you to use as a reference for developing your own plugins.

[![Main](https://github.com/SergioFloresG/pantrypath/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/SergioFloresG/traefik-cors-middleware/actions/workflows/main.yml)
[![Matrix](https://github.com/SergioFloresG/pantrypath/actions/workflows/go-cross.yml/badge.svg?branch=master)](https://github.com/SergioFloresG/traefik-cors-middleware/actions/workflows/go-cross.yml)

The existing plugins can be browsed into the [Plugin Catalog](https://plugins.traefik.io).

## Usage

Use this middleware to add CORS support to the necessary endpoints, indicating their domain, allowed methods and more.

### Configuration

Requirements: Traefik >= v2.5.5

| Option              | Description                                                                                                                                         | Header                                                                                                                         |
|---------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------|
| `allow_credentials` | This indicates whether or not the actual request can be made using credentials.<br/>_(default: `false`)_                                            | [Access Control Allow Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials) | 
| `allow_origins`     | Indicates the group of sources that have access to the requesting resource.<br/>_(default: `*`)_                                                    | [Access Control Allow Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin)           |
| `allow_methods`     | Specifies one or more methods allowed when accessing a resource in response.<br/>_(default: `OPTIONS`, `GET`, `POST`, `PUT`, `DELETE`)_             | [Access Control Allow Methods](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods)         |
| `allow_headers`     | Indicate which HTTP headers can be used during the actual request. the indicated headers are an addition to [default ones](#allow_headers_defaults) | [Access Control Allow Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers)         |
| `expose_headers`    | Indicate which response headers should be made available to scripts running in the browser, in response to a cross-origin request.                  | [Access Control Expose Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers)       |
| `max_age`:          | Indicates how long the results of a preflight request can be cached <br/>_(default: `86400`)_                                                      | [Access Control Max Age](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age)                     |

<a name="allow_headers_defaults"></a>
**Default Allow Headers**: Content-Type, Content-Length, Accept-Encoding, Authorization, Accept, Origin, Referer,
Cache-Control

#### Traefik

```yaml
# Static configuration

experimental:
  plugins:
    example:
      moduleName: github.com/SergioFloresG/traefik-cors-middleware
      version: v0.1.0
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is
the `http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-cors-middleware

  services:
    service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    my-cors-middleware:
      plugin:
        corsmiddleware:
          allow_origins: ["https://example.com"]
```