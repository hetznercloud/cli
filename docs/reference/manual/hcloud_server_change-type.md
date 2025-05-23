## hcloud server change-type

Change type of a server

```
hcloud server change-type [--keep-disk] <server> <server-type>
```

### Options

```
  -h, --help        help for change-type
      --keep-disk   Keep disk size of current Server Type. This enables downgrading the server. (true, false)
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
