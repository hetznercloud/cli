## hcloud volume create

Create a Volume

```
hcloud volume create [options] --name <name> --size <size>
```

### Options

```
      --automount[=true|false]      Automount Volume after attach (Server must be provided)
      --enable-protection strings   Enable protection (delete) (default: none)
      --format string               Format Volume after creation (ext4 or xfs)
  -h, --help[=true|false]           help for create
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string             Location (ID or name)
      --name string                 Volume name (required)
  -o, --output stringArray          output options: json|yaml
      --server string               Server (ID or name)
      --size int                    Size (GB) (required)
```

### Options inherited from parent commands

```
      --config string                           Config file path (default "~/.config/hcloud/cli.toml")
      --context string                          Currently active context
      --debug[=true|false]                      Enable debug output
      --debug-file string                       File to write debug output to
      --endpoint string                         Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --hetzner-endpoint string                 Hetzner API endpoint (default "https://api.hetzner.com/v1")
      --http-timeout duration                   Timeout for HTTP requests (default 0 = no timeout)
      --no-experimental-warnings[=true|false]   If true, experimental warnings are not shown
      --poll-interval duration                  Interval at which to poll information, for example action progress (default 500ms)
      --quiet[=true|false]                      If true, only print error messages
```

### SEE ALSO

* [hcloud volume](hcloud_volume.md)	 - Manage Volumes

