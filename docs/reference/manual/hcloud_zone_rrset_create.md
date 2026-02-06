## hcloud zone rrset create

Create a Zone RRSet

### Synopsis

Create a Zone RRSet.

The optional records file has to be in JSON format. You can find the schema at https://docs.hetzner.cloud/reference/cloud#zone-rrset-actions-set-records-of-an-rrset

Example file content:

[
  {
    "value": "198.51.100.1",
    "comment": "My web server at Hetzner Cloud."
  },
  {
    "value": "198.51.100.2",
    "comment": "My other server at Hetzner Cloud."
  }
]

```
hcloud zone rrset create [options] --name <name> --type <A|AAAA|CAA|CNAME|DS|HINFO|HTTPS|MX|NS|PTR|RP|SOA|SRV|SVCB|TLSA|TXT> (--record <record>... | --records-file <file>) <zone>
```

### Options

```
  -h, --help                   help for create
      --label stringToString   User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string            Name of the Zone RRSet (required)
  -o, --output stringArray     output options: json|yaml
      --record stringArray     List of records (can be specified multiple times, conflicts with --records-file)
      --records-file string    JSON file containing the records (conflicts with --record)
      --ttl int                Time To Live (TTL) of the RRSet
      --type string            Type of the Zone RRSet (required)
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

* [hcloud zone rrset](hcloud_zone_rrset.md)	 - Manage Zone RRSets (records)

