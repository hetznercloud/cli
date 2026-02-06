## hcloud storage-box subaccount update-access-settings

Update access settings of the Storage Box Subaccount

```
hcloud storage-box subaccount update-access-settings [options] <storage-box> <subaccount>
```

### Options

```
      --enable-samba           Whether the Samba subsystem should be enabled (true, false)
      --enable-ssh             Whether the SSH subsystem should be enabled (true, false)
      --enable-webdav          Whether the WebDAV subsystem should be enabled (true, false)
  -h, --help                   help for update-access-settings
      --reachable-externally   Whether the Storage Box should be accessible from outside the Hetzner network (true, false)
      --readonly               Whether the Subaccount should be read-only (true, false)
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

* [hcloud storage-box subaccount](hcloud_storage-box_subaccount.md)	 - Manage Storage Box Subaccounts

