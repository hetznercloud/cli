## hcloud server ssh

Spawn an SSH connection for the Server

```
hcloud server ssh [options] <server> [--] [ssh options] [command [argument...]]
```

### Options

```
  -h, --help[=true|false]   help for ssh
      --ipv6[=true|false]   Establish SSH connection to IPv6 address
  -p, --port int            Port for SSH connection (default 22)
  -u, --user string         Username for SSH connection (default "root")
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

* [hcloud server](hcloud_server.md)	 - Manage Servers

