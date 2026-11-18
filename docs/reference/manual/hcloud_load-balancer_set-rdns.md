## hcloud load-balancer set-rdns

Change reverse DNS of a Load Balancer

```
hcloud load-balancer set-rdns [--ip <ip>] (--hostname <hostname> | --reset) <load-balancer>
```

### Options

```
  -h, --help[=true|false]    help for set-rdns
  -r, --hostname string      Hostname to set as a reverse DNS PTR entry
  -i, --ip ip                IP address for which the reverse DNS entry should be set
      --reset[=true|false]   Reset the reverse DNS entry to the default value
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

* [hcloud load-balancer](hcloud_load-balancer.md)	 - Manage Load Balancers

