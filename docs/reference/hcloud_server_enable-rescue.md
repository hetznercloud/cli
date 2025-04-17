## hcloud server enable-rescue

Enable rescue for a server

```
hcloud server enable-rescue [options] <server>
```

### Options

```
  -h, --help              help for enable-rescue
      --ssh-key strings   ID or name of SSH Key to inject (can be specified multiple times)
      --type string       Rescue type (default "linux64")
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
