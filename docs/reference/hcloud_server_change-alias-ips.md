## hcloud server change-alias-ips

Change a server's alias IPs in a Network

```
hcloud server change-alias-ips [options] --network <network> <server>
```

### Options

```
      --alias-ips strings   New alias IPs
      --clear               Remove all alias IPs
  -h, --help                help for change-alias-ips
  -n, --network string      Network (ID or name) (required)
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
