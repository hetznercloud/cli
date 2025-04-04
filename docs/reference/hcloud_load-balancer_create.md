## hcloud load-balancer create

Create a Load Balancer

```
hcloud load-balancer create [options] --name <name> --type <type>
```

### Options

```
      --algorithm-type string       Algorithm Type name (round_robin or least_connections)
      --enable-protection strings   Enable protection (delete) (default: none)
  -h, --help                        help for create
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string             Location (ID or name)
      --name string                 Load Balancer name (required)
      --network string              Name or ID of the Network the Load Balancer should be attached to on creation
      --network-zone string         Network Zone
  -o, --output stringArray          output options: json|yaml
      --type string                 Load Balancer Type (ID or name) (required)
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
