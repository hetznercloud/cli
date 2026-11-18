## hcloud network expose-routes-to-vswitch

Expose routes to connected vSwitch

### Synopsis

Enabling this will expose routes to the connected vSwitch. Set the --disable flag to remove the exposed routes.

```
hcloud network expose-routes-to-vswitch [--disable] <network>
```

### Options

```
      --disable[=true|false]   Remove any exposed routes from the connected vSwitch
  -h, --help[=true|false]      help for expose-routes-to-vswitch
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

* [hcloud network](hcloud_network.md)	 - Manage Networks

