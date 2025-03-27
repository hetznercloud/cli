## hcloud load-balancer remove-target

Remove a target from a Load Balancer

```
hcloud load-balancer remove-target [options] <load-balancer>
```

### Options

```
  -h, --help                    help for remove-target
      --ip string               IP address of an IP target
      --label-selector string   Label Selector
      --server string           Name or ID of the server
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
