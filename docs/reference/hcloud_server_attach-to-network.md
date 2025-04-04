## hcloud server attach-to-network

Attach a Server to a Network

```
hcloud server attach-to-network [options] --network <network> <server>
```

### Options

```
      --alias-ips ipSlice   Additional IP addresses to be assigned to the Server (default [])
  -h, --help                help for attach-to-network
      --ip ip               IP address to assign to the Server (auto-assigned if omitted)
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
