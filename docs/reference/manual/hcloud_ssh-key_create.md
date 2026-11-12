## hcloud ssh-key create

Create an SSH Key

```
hcloud ssh-key create [options] --name <name> (--public-key <key> | --public-key-from-file <file>)
```

### Options

```
  -h, --help[=true|false]             help for create
      --label stringToString          User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string                   Key name (required)
  -o, --output stringArray            output options: json|yaml
      --public-key string             Public key
      --public-key-from-file string   Path to file containing public key
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

* [hcloud ssh-key](hcloud_ssh-key.md)	 - Manage SSH Keys

