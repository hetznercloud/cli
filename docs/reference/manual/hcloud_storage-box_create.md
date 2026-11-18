## hcloud storage-box create

Create a new Storage Box

```
hcloud storage-box create [options] --name <name> --type <type> --location <location> --password <password>
```

### Options

```
      --enable-protection strings           Enable protection (delete) (default: none)
      --enable-samba[=true|false]           Whether the Samba subsystem should be enabled
      --enable-ssh[=true|false]             Whether the SSH subsystem should be enabled
      --enable-webdav[=true|false]          Whether the WebDAV subsystem should be enabled
      --enable-zfs[=true|false]             Whether the ZFS Snapshot folder should be visible
  -h, --help[=true|false]                   help for create
      --label stringToString                User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string                     Location (ID or name) (required)
      --name string                         Storage Box name (required)
  -o, --output stringArray                  output options: json|yaml
      --password string                     The password that will be set for this Storage Box (required)
      --reachable-externally[=true|false]   Whether the Storage Box should be accessible from outside the Hetzner network
      --ssh-key stringArray                 SSH public keys in OpenSSH format or as the ID or name of an existing SSH key
      --type string                         Storage Box Type (ID or name) (required)
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

* [hcloud storage-box](hcloud_storage-box.md)	 - Manage Storage Boxes

