## hcloud server metrics

[ALPHA] Metrics from a Server

```
hcloud server metrics [options] (--type <cpu|disk|network>)... <server>
```

### Options

```
      --end string           ISO 8601 timestamp
  -h, --help                 help for metrics
  -o, --output stringArray   output options: json|yaml
      --start string         ISO 8601 timestamp
      --type strings         Types of metrics you want to show
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
