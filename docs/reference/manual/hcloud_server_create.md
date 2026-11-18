## hcloud server create

Create a Server

### Synopsis

Create a Server.

The --datacenter flag has been removed. Use --location instead.
See https://docs.hetzner.cloud/changelog#2026-07-01-removing-datacenters.

```
hcloud server create [options] --name <name> --type <server-type> --image <image>
```

### Options

```
      --allow-deprecated-image[=true|false]   Enable the use of deprecated Images
      --automount[=true|false]                Automount Volumes after attach
      --enable-backup[=true|false]            Enable automatic backups
      --enable-protection strings             Enable protection (delete, rebuild) (default: none)
      --firewall strings                      ID or name of Firewall to attach the Server to (can be specified multiple times)
  -h, --help[=true|false]                     help for create
      --image string                          Image (ID or name) (required)
      --label stringToString                  User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string                       Location (ID or name)
      --name string                           Server name (required)
      --network strings                       ID or name of Network to attach the Server to (can be specified multiple times)
  -o, --output stringArray                    output options: json|yaml
      --placement-group string                Placement Group (ID of name)
      --primary-ipv4 string                   Primary IPv4 (ID of name)
      --primary-ipv6 string                   Primary IPv6 (ID of name)
      --ssh-key strings                       ID or name of SSH Key to inject (can be specified multiple times)
      --start-after-create[=true|false]       Start Server right after creation (default true)
      --type string                           Server Type (ID or name) (required)
      --user-data-from-file stringArray       Read user data from specified file (use - to read from stdin)
      --volume strings                        ID or name of Volume to attach (can be specified multiple times)
      --without-ipv4[=true|false]             Creates the Server without an IPv4
      --without-ipv6[=true|false]             Creates the Server without an IPv6
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

* [hcloud server](hcloud_server.md)	 - Manage Servers

