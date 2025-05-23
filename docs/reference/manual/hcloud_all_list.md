## hcloud all list

List all resources in the project

### Synopsis

List all resources in the project. This does not include static/public resources like Locations, public ISOs, etc.

Listed resources are:
 - Servers
 - Images
 - Placement Groups
 - Primary IPs
 - ISOs
 - Volumes
 - Load Balancer
 - Floating IPs
 - Networks
 - Firewalls
 - Certificates
 - SSH Keys

```
hcloud all list [options]
```

### Options

```
  -h, --help                 help for list
  -o, --output stringArray   output options: json|yaml
      --paid                 Only list resources that cost money (true, false)
  -l, --selector string      Selector to filter by labels
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

* [hcloud all](hcloud_all.md)	 - Commands that apply to all resources
