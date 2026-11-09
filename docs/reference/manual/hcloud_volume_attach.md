## hcloud volume attach

Attach a Volume to a Server

### Synopsis

Attach a Volume to a Server.

The --server flag is deprecated. Use --assignee and --assignee-type instead.
See https://docs.hetzner.cloud/changelog#2026-10-27-volumes-allow-generic-resource-targets.

```
hcloud volume attach [--automount] --assignee-type <type> --assignee <assignee> <volume>
```

### Options

```
      --assignee string        Assignee (ID, or name if the assignee type is server)
      --assignee-type string   Assignee type (server)
      --automount              Automount Volume after attach (true, false)
  -h, --help                   help for attach
      --server string          Server (ID or name) (deprecated)
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

