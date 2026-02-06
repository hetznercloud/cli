## hcloud zone import-zonefile

Imports a zone file, replacing all Zone RRSets

```
hcloud zone import-zonefile --zonefile <file> <zone>
```

### Options

```
  -h, --help                 help for import-zonefile
  -o, --output stringArray   output options: json|yaml
      --zonefile string      Zone file in BIND (RFC 1034/1035) format (use - to read from stdin)
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

