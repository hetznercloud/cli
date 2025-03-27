## hcloud network remove-route

Remove a route from a network

```
hcloud network remove-route --destination <destination> --gateway <ip> <network>
```

### Options

```
      --destination ipNet   Destination network or host (required)
      --gateway ip          Gateway IP address (required)
  -h, --help                help for remove-route
```

### Options inherited from parent commands

```
      --config string            Config file path (default "/Users/paul/.config/hcloud/cli.toml")
      --context string           Currently active context
      --debug                    Enable debug output
      --debug-file string        File to write debug output to
      --endpoint string          Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --poll-interval duration   Interval at which to poll information, for example action progress (default 500ms)
      --quiet                    If true, only print error messages
```

### SEE ALSO

* [hcloud network](hcloud_network.md)	 - Manage networks
