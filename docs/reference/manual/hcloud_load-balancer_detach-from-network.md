## hcloud load-balancer detach-from-network

Detach a Load Balancer from a Network

```
hcloud load-balancer detach-from-network --network <network> <load-balancer>
```

### Options

```
  -h, --help             help for detach-from-network
  -n, --network string   Network (ID or name) (required)
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
