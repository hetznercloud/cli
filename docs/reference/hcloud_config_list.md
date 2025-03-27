## hcloud config list

List configuration values

```
hcloud config list
```

### Options

```
  -a, --all                  Also show default values
      --allow-sensitive      Allow showing sensitive values
  -g, --global               Only show global values
  -h, --help                 help for list
  -o, --output stringArray   output options: noheader|columns=...|json|yaml
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

* [hcloud config](hcloud_config.md)	 - Manage configuration
