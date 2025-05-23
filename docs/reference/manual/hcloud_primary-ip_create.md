## hcloud primary-ip create

Create a Primary IP

```
hcloud primary-ip create [options] --type <ipv4|ipv6> --name <name>
```

### Options

```
      --assignee-id int             Assignee (usually a Server) to assign Primary IP to
      --auto-delete                 Delete Primary IP if assigned resource is deleted (true, false)
      --datacenter string           Datacenter (ID or name)
      --enable-protection strings   Enable protection (delete) (default: none)
  -h, --help                        help for create
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string                 Name (required)
  -o, --output stringArray          output options: json|yaml
      --type string                 Type (ipv4 or ipv6) (required)
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

* [hcloud primary-ip](hcloud_primary-ip.md)	 - Manage Primary IPs
