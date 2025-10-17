## hcloud zone rrset add-label

[experimental] Add a label to a Zone RRSet

### Synopsis

Add a label to a Zone RRSet

Experimental: DNS API is in beta, breaking changes may occur within minor releases.
See https://docs.hetzner.cloud/changelog#2025-10-07-dns-beta for more details.


```
hcloud zone rrset add-label [--overwrite] <zone> <name> <type> <label>...
```

### Options

```
  -h, --help        help for add-label
  -o, --overwrite   Overwrite label if it exists already (true, false)
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

* [hcloud zone rrset](hcloud_zone_rrset.md)	 - [experimental] Manage Zone RRSets (records)
