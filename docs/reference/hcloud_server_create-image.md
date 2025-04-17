## hcloud server create-image

Create an Image from a Server

```
hcloud server create-image [options] --type <snapshot|backup> <server>
```

### Options

```
      --description string     Image description
  -h, --help                   help for create-image
      --label stringToString   User-defined labels ('key=value') (can be specified multiple times) (default [])
      --type string            Image type (required)
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
