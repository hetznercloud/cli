## hcloud placement-group create

Create a Placement Group

```
hcloud placement-group create [options] --name <name> --type <type>
```

### Options

```
  -h, --help                   help for create
      --label stringToString   User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string            Name
  -o, --output stringArray     output options: json|yaml
      --type string            Type of the Placement Group
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

* [hcloud placement-group](hcloud_placement-group.md)	 - Manage Placement Groups
