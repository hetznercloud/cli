## hcloud config remove

Remove values from a list

### Synopsis

Remove values from a list. For a list of all available configuration options, run 'hcloud help config'.

```
hcloud config remove <key> <value>...
```

### Options

```
      --global   Remove the value(s) globally (for all contexts) (true, false)
  -h, --help     help for remove
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

* [hcloud config](hcloud_config.md)	 - Manage configuration
