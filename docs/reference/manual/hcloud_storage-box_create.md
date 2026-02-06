## hcloud storage-box create

Create a new Storage Box

```
hcloud storage-box create [options] --name <name> --type <type> --location <location> --password <password>
```

### Options

```
      --enable-protection strings   Enable protection (delete) (default: none)
      --enable-samba                Whether the Samba subsystem should be enabled (true, false)
      --enable-ssh                  Whether the SSH subsystem should be enabled (true, false)
      --enable-webdav               Whether the WebDAV subsystem should be enabled (true, false)
      --enable-zfs                  Whether the ZFS Snapshot folder should be visible (true, false)
  -h, --help                        help for create
      --label stringToString        User-defined labels ('key=value') (can be specified multiple times) (default [])
      --location string             Location (ID or name) (required)
      --name string                 Storage Box name (required)
  -o, --output stringArray          output options: json|yaml
      --password string             The password that will be set for this Storage Box (required)
      --reachable-externally        Whether the Storage Box should be accessible from outside the Hetzner network (true, false)
      --ssh-key stringArray         SSH public keys in OpenSSH format or as the ID or name of an existing SSH key
      --type string                 Storage Box Type (ID or name) (required)
```

### Options inherited from parent commands

```
      --config string              Config file path (default "~/.config/hcloud/cli.toml")
      --context string             Currently active context
      --debug                      Enable debug output
      --debug-file string          File to write debug output to
      --endpoint string            Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --hetzner-endpoint string    Hetzner API endpoint (default "https://api.hetzner.com/v1")
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud storage-box](hcloud_storage-box.md)	 - Manage Storage Boxes

