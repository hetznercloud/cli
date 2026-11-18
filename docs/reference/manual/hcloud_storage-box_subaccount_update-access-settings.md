## hcloud storage-box subaccount update-access-settings

Update access settings of the Storage Box Subaccount

```
hcloud storage-box subaccount update-access-settings [options] <storage-box> <subaccount>
```

### Options

```
      --enable-samba[=true|false]           Whether the Samba subsystem should be enabled
      --enable-ssh[=true|false]             Whether the SSH subsystem should be enabled
      --enable-webdav[=true|false]          Whether the WebDAV subsystem should be enabled
  -h, --help[=true|false]                   help for update-access-settings
      --reachable-externally[=true|false]   Whether the Storage Box should be accessible from outside the Hetzner network
      --readonly[=true|false]               Whether the Subaccount should be read-only
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

* [hcloud storage-box subaccount](hcloud_storage-box_subaccount.md)	 - Manage Storage Box Subaccounts

