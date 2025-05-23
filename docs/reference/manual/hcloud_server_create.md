## hcloud server create

Create a Server

```
hcloud server create [options] --name <name> --type <server-type> --image <image>
```

### Options

```
      --allow-deprecated-image            Enable the use of deprecated Images (default: false) (true, false)
      --automount                         Automount Volumes after attach (default: false) (true, false)
      --datacenter string                 Datacenter (ID or name)
      --enable-backup                     Enable automatic backups (true, false)
      --enable-protection strings         Enable protection (delete, rebuild) (default: none)
      --firewall strings                  ID or name of Firewall to attach the Server to (can be specified multiple times)
  -h, --help                              help for create
      --image string                      Image (ID or name) (required)
      --label stringToString              User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string                   Location (ID or name)
      --name string                       Server name (required)
      --network strings                   ID or name of Network to attach the Server to (can be specified multiple times)
  -o, --output stringArray                output options: json|yaml
      --placement-group string            Placement Group (ID of name)
      --primary-ipv4 string               Primary IPv4 (ID of name)
      --primary-ipv6 string               Primary IPv6 (ID of name)
      --ssh-key strings                   ID or name of SSH Key to inject (can be specified multiple times)
      --start-after-create                Start Server right after creation (true, false) (default true)
      --type string                       Server Type (ID or name) (required)
      --user-data-from-file stringArray   Read user data from specified file (use - to read from stdin)
      --volume strings                    ID or name of Volume to attach (can be specified multiple times)
      --without-ipv4                      Creates the Server without an IPv4 (default: false) (true, false)
      --without-ipv6                      Creates the Server without an IPv6 (default: false) (true, false)
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
