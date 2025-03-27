## hcloud load-balancer list

List Load Balancer

### Synopsis

Displays a list of Load Balancer.

Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=age,created (see available columns below).

Columns:
 - age
 - created
 - health
 - id
 - ipv4
 - ipv6
 - labels
 - location
 - name
 - network_zone
 - protection
 - type

```
hcloud load-balancer list [options]
```

### Options

```
  -h, --help                 help for list
  -o, --output stringArray   output options: noheader|columns=...|json|yaml
  -l, --selector string      Selector to filter by labels
  -s, --sort strings         Determine the sorting of the result
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

* [hcloud load-balancer](hcloud_load-balancer.md)	 - Manage Load Balancers
