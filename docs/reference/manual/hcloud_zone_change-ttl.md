## hcloud zone change-ttl

Changes the default Time To Live (TTL) of a Zone

```
hcloud zone change-ttl --ttl <ttl> <zone>
```

### Options

```
  -h, --help      help for change-ttl
      --ttl int   Default Time To Live (TTL) of the Zone (required) (default 3600)
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

* [hcloud zone](hcloud_zone.md)	 - Manage DNS Zones and Zone RRSets (records)

