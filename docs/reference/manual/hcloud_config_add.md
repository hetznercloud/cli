## hcloud config add

Add values to a list

### Synopsis

Add values to a list. For a list of all available configuration options, run 'hcloud help config'.

```
hcloud config add <key> <value>...
```

### Options

```
      --global[=true|false]   Set the value globally (for all contexts)
  -h, --help[=true|false]     help for add
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

* [hcloud config](hcloud_config.md)	 - Manage configuration

