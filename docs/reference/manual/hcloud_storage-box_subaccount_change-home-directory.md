## hcloud storage-box subaccount change-home-directory

Update access settings of the Storage Box Subaccount

```
hcloud storage-box subaccount change-home-directory --home-directory <home-directory> <storage-box> <subaccount>
```

### Options

```
  -h, --help                    help for change-home-directory
      --home-directory string   Home directory of the Subaccount. Will be created if it doesn't exist yet
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

