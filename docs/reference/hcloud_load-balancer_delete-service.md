## hcloud load-balancer delete-service

Deletes a service from a Load Balancer

```
hcloud load-balancer delete-service --listen-port <1-65535> <load-balancer>
```

### Options

```
  -h, --help              help for delete-service
      --listen-port int   The listen port of the service you want to delete (required)
```

### Options inherited from parent commands

```
      --config string            Config file path (default "~/.config/hcloud/cli.toml")
      --context string           Currently active context
      --debug                    Enable debug output
      --debug-file string        File to write debug output to
      --endpoint string          Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --poll-interval duration   Interval at which to poll information, for example action progress (default 500ms)
      --quiet                    If true, only print error messages
```

### SEE ALSO

* [hcloud load-balancer](hcloud_load-balancer.md)	 - Manage Load Balancers
