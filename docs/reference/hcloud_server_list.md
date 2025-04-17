## hcloud server list

List Servers

### Synopsis

Displays a list of Servers.

Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=age,backup_window (see available columns below).

Columns:
 - age
 - backup_window
 - created
 - datacenter
 - id
 - included_traffic
 - ingoing_traffic
 - ipv4
 - ipv6
 - labels
 - location
 - locked
 - name
 - outgoing_traffic
 - placement_group
 - primary_disk_size
 - private_net
 - protection
 - rescue_enabled
 - status
 - type
 - volumes

```
hcloud server list [options]
```

### Options

```
  -h, --help                 help for list
  -o, --output stringArray   output options: noheader|columns=...|json|yaml
  -l, --selector string      Selector to filter by labels
  -s, --sort strings         Determine the sorting of the result
      --status strings       Only Servers with one of these statuses are displayed
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
