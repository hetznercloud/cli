## hcloud storage-box change-type

[experimental] Change type of a Storage Box

### Synopsis

Requests a Storage Box to be upgraded or downgraded to another Storage Box Type.
Please note that it is not possible to downgrade to a Storage Box Type that offers less disk space than you are currently using.

Experimental: Storage Box support is experimental, breaking changes may occur within minor releases.
See https://github.com/hetznercloud/cli/issues/1202 for more details.


```
hcloud storage-box change-type <storage-box> <storage-box-type>
```

### Options

```
  -h, --help   help for change-type
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
