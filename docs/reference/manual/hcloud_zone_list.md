## hcloud zone list

[experimental] List Zones

### Synopsis

Displays a list of Zones.

Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=age,authoritative_nameservers (see available columns below).

Columns:
 - age
 - authoritative_nameservers
 - created
 - id
 - labels
 - mode
 - name
 - name_idna
 - primary_nameservers
 - protection
 - record_count
 - registrar
 - status
 - ttl

Experimental: DNS API is in beta, breaking changes may occur within minor releases.
See https://docs.hetzner.cloud/changelog#2025-10-07-dns-beta for more details.


```
hcloud zone list [options]
```

### Options

```
  -h, --help                 help for list
      --mode string          Only Zones with this mode are displayed
  -o, --output stringArray   output options: noheader|columns=...|json|yaml
  -l, --selector string      Selector to filter by labels
  -s, --sort strings         Determine the sorting of the result
```

### Options inherited from parent commands

```
      --config string              Config file path (default "~/.config/hcloud/cli.toml")
      --context string             Currently active context
      --debug                      Enable debug output
      --debug-file string          File to write debug output to
      --endpoint string            Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud zone](hcloud_zone.md)	 - [experimental] Manage DNS Zones and Zone RRSets (records)
