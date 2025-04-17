## hcloud firewall remove-from-resource

Removes a Firewall from a single resource

```
hcloud firewall remove-from-resource (--type server --server <server> | --type label_selector --label-selector <label-selector>) <firewall>
```

### Options

```
  -h, --help                    help for remove-from-resource
  -l, --label-selector string   Label Selector
      --server string           Server name of ID (required when type is server)
      --type string             Resource Type (server) (required)
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

* [hcloud firewall](hcloud_firewall.md)	 - Manage Firewalls
