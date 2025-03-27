## hcloud firewall create

Create a Firewall

```
hcloud firewall create [options] --name <name>
```

### Options

```
  -h, --help                   help for create
      --label stringToString   User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string            Name
  -o, --output stringArray     output options: json|yaml
      --rules-file string      JSON file containing your routes (use - to read from stdin). The structure of the file needs to be the same as within the API: https://docs.hetzner.cloud/#firewalls-get-a-firewall 
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

* [hcloud firewall](hcloud_firewall.md)	 - Manage Firewalls
