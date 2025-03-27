## hcloud volume create

Create a volume

```
hcloud volume create [options] --name <name> --size <size>
```

### Options

```
      --automount                   Automount volume after attach (server must be provided)
      --enable-protection strings   Enable protection (delete) (default: none)
      --format string               Format volume after creation (ext4 or xfs)
  -h, --help                        help for create
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string             Location (ID or name)
      --name string                 Volume name (required)
  -o, --output stringArray          output options: json|yaml
      --server string               Server (ID or name)
      --size int                    Size (GB) (required)
```

### Options inherited from parent commands

```
      --config string            Config file path (default "/Users/paul/.config/hcloud/cli.toml")
      --context string           Currently active context
      --debug                    Enable debug output
      --debug-file string        File to write debug output to
      --endpoint string          Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --poll-interval duration   Interval at which to poll information, for example action progress (default 500ms)
      --quiet                    If true, only print error messages
```

### SEE ALSO

* [hcloud volume](hcloud_volume.md)	 - Manage Volumes
