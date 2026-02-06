## hcloud zone change-primary-nameservers

Changes the primary nameservers of a secondary Zone

### Synopsis

Changes the primary nameservers of a secondary Zone.

Input file has to be in JSON format. You can find the schema at
https://docs.hetzner.cloud/reference/cloud#zone-actions-change-a-zones-primary-nameservers

Example file content:

[
  {
    "address": "203.0.113.10"
  },
  {
    "address": "203.0.113.11",
    "port": 5353
  },
  {
    "address": "203.0.113.12",
    "tsig_algorithm": "hmac-sha256",
    "tsig_key": "example-key"
  }
]

```
hcloud zone change-primary-nameservers --primary-nameservers-file <file> <zone>
```

### Options

```
  -h, --help                              help for change-primary-nameservers
      --primary-nameservers-file string   JSON file containing the new primary nameservers. (use - to read from stdin)
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

