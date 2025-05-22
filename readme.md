# headers-extra

[Traefik](https://traefik.io) plugins are developed using [http-wasm](https://http-wasm.io/).

## Usage

see configuration


### Configuration

For each plugin, the Traefik static configuration must define the module name (as is usual for Go packages).

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration

experimental:
  plugins:
    headers_extra:
      moduleName: github.com/gfx-labs/plugin-traefik-headers-extra
      version: v0.0.1
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

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
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  middlewares:
    my-plugin:
      plugin:
        headers_extra:
          copy:
            - from: X-Real-IP
              to: X-Icarus-Real-IP
```
