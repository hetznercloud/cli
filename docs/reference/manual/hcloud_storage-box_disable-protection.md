## hcloud storage-box disable-protection

[experimental] Disable resource protection for a Storage Box

### Synopsis

Disable resource protection for a Storage Box

Experimental: Storage Box support is experimental, breaking changes may occur within minor releases.
See https://github.com/hetznercloud/cli/issues/1202 for more details.


```
hcloud storage-box disable-protection <storage-box> delete
```

### Options

```
  -h, --help   help for disable-protection
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

* [hcloud storage-box](hcloud_storage-box.md)	 - [experimental] Manage Storage Boxes
