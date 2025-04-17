## hcloud network add-subnet

Add a subnet to a Network

```
hcloud network add-subnet [options] --type <cloud|server|vswitch> --network-zone <zone> <network>
```

### Options

```
  -h, --help                  help for add-subnet
      --ip-range ipNet        Range to allocate IPs from
      --network-zone string   Name of Network zone (required)
      --type string           Type of subnet (required)
      --vswitch-id int        ID of the vSwitch
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

* [hcloud network](hcloud_network.md)	 - Manage Networks
