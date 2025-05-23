## hcloud load-balancer add-target

Add a target to a Load Balancer

```
hcloud load-balancer add-target [options] (--server <server> | --label-selector <label-selector> | --ip <ip>) <load-balancer>
```

### Options

```
  -h, --help                    help for add-target
      --ip string               Use the passed IP address as target
      --label-selector string   Label Selector
      --server string           Name or ID of the server
      --use-private-ip          Determine if the Load Balancer should connect to the target via the network (true, false)
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
