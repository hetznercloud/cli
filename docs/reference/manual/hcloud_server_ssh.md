## hcloud server ssh

Spawn an SSH connection for the Server

```
hcloud server ssh [options] <server> [--] [ssh options] [command [argument...]]
```

### Options

```
  -h, --help          help for ssh
      --ipv6          Establish SSH connection to IPv6 address (true, false)
  -p, --port int      Port for SSH connection (default 22)
  -u, --user string   Username for SSH connection (default "root")
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

* [hcloud server](hcloud_server.md)	 - Manage Servers
