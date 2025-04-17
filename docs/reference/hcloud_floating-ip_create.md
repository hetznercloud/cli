## hcloud floating-ip create

Create a Floating IP

```
hcloud floating-ip create [options] --type <ipv4|ipv6> (--home-location <location> | --server <server>)
```

### Options

```
      --description string          Description
      --enable-protection strings   Enable protection (delete) (default: none)
  -h, --help                        help for create
      --home-location string        Home Location
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string                 Name
  -o, --output stringArray          output options: json|yaml
      --server string               Server to assign Floating IP to
      --type string                 Type (ipv4 or ipv6) (required)
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

* [hcloud floating-ip](hcloud_floating-ip.md)	 - Manage Floating IPs
