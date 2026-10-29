## hcloud volume create

Create a Volume

### Synopsis

Create a Volume.

The --server flag is deprecated. Use --assignee and --assignee-type instead.
See https://docs.hetzner.cloud/changelog#2026-10-27-volumes-allow-generic-resource-targets.

```
hcloud volume create [options] --name <name> --size <size>
```

### Options

```
      --assignee string             Assignee (ID, or name if the assignee type is server)
      --assignee-type string        Assignee type (server)
      --automount                   Automount Volume after attach (assignee must be provided) (true, false)
      --enable-protection strings   Enable protection (delete) (default: none)
      --format string               Format Volume after creation (ext4 or xfs)
  -h, --help                        help for create
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string             Location (ID or name)
      --name string                 Volume name (required)
  -o, --output stringArray          output options: json|yaml
      --server string               Server (ID or name) (deprecated)
      --size int                    Size (GB) (required)
```

### Options inherited from parent commands

```
      --config string              Config file path (default "~/.config/hcloud/cli.toml")
      --context string             Currently active context
      --debug                      Enable debug output
      --debug-file string          File to write debug output to
      --endpoint string            Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --hetzner-endpoint string    Hetzner API endpoint (default "https://api.hetzner.com/v1")
      --http-timeout duration      Timeout for HTTP requests (default 0 = no timeout)
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud volume](hcloud_volume.md)	 - Manage Volumes

