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
 - Storage Boxes
 - Zones

```
hcloud all list [options]
```

### Options

```
  -h, --help[=true|false]    help for list
  -o, --output stringArray   output options: json|yaml
      --paid[=true|false]    Only list resources that cost money
  -l, --selector string      Selector to filter by labels
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

* [hcloud all](hcloud_all.md)	 - Commands that apply to all resources

