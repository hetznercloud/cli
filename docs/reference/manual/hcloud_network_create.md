## hcloud network create

Create a Network

```
hcloud network create [options] --name <name> --ip-range <ip-range>
```

### Options

```
      --enable-protection strings   Enable protection (delete) (default: none)
      --expose-routes-to-vswitch    Expose routes from this Network to the vSwitch connection. It only takes effect if a vSwitch connection is active. (true, false)
  -h, --help                        help for create
      --ip-range ipNet              Network IP range (required)
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string                 Network name (required)
  -o, --output stringArray          output options: json|yaml
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

* [hcloud network](hcloud_network.md)	 - Manage Networks
