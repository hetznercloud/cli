## hcloud storage-box subaccount update

[experimental] Update a Storage Box Subaccount

### Synopsis

Update a Storage Box Subaccount

Experimental: Storage Box support is experimental, breaking changes may occur within minor releases.
See https://github.com/hetznercloud/cli/issues/1202 for more details.


```
hcloud storage-box subaccount update [options] <storage-box> <subaccount>
```

### Options

```
      --description string   Description of the Storage Box Subaccount
  -h, --help                 help for update
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

* [hcloud storage-box subaccount](hcloud_storage-box_subaccount.md)	 - [experimental] Manage Storage Box Subaccounts
