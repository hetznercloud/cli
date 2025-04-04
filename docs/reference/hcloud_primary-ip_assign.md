## hcloud primary-ip assign

Assign a Primary IP to an assignee (usually a Server)

```
hcloud primary-ip assign --server <server> <primary-ip>
```

### Options

```
  -h, --help            help for assign
      --server string   Name or ID of the Server
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

* [hcloud primary-ip](hcloud_primary-ip.md)	 - Manage Primary IPs
