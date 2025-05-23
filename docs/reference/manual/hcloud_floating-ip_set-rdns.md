## hcloud floating-ip set-rdns

Change reverse DNS of a Floating IP

```
hcloud floating-ip set-rdns [--ip <ip>] (--hostname <hostname> | --reset) <floating-ip>
```

### Options

```
  -h, --help              help for set-rdns
  -r, --hostname string   Hostname to set as a reverse DNS PTR entry
  -i, --ip ip             IP address for which the reverse DNS entry should be set
      --reset             Reset the reverse DNS entry to the default value (true, false)
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

* [hcloud floating-ip](hcloud_floating-ip.md)	 - Manage Floating IPs
