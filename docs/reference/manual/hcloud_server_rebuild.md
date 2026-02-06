## hcloud server rebuild

Rebuild a server

```
hcloud server rebuild [options] --image <image> <server>
```

### Options

```
      --allow-deprecated-image            Enable the use of deprecated images (default: false) (true, false)
  -h, --help                              help for rebuild
      --image string                      ID or name of Image to rebuild from (required)
      --user-data-from-file stringArray   Read user data from specified file (use - to read from stdin)
```

### Options inherited from parent commands

```
      --config string              Config file path (default "~/.config/hcloud/cli.toml")
      --context string             Currently active context
      --debug                      Enable debug output
      --debug-file string          File to write debug output to
      --endpoint string            Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --hetzner-endpoint string    Hetzner API endpoint (default "https://api.hetzner.com/v1")
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud server](hcloud_server.md)	 - Manage Servers

