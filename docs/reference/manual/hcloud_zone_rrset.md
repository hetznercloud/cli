## hcloud zone rrset

Manage Zone RRSets (records)

### Synopsis


For more details, see the documentation for Zones https://docs.hetzner.cloud/reference/cloud#zones
or Zone RRSets https://docs.hetzner.cloud/reference/cloud#zone-rrsets.

TXT records format must consist of one or many quoted strings of 255 characters. If the
user provider TXT records are not quoted, they will be formatted for you.

### Options

```
  -h, --help   help for rrset
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
* [hcloud zone rrset add-label](hcloud_zone_rrset_add-label.md)	 - Add a label to a Zone RRSet
* [hcloud zone rrset add-records](hcloud_zone_rrset_add-records.md)	 - Add records to a Zone RRSet
* [hcloud zone rrset change-ttl](hcloud_zone_rrset_change-ttl.md)	 - Changes the Time To Live (TTL) of a Zone RRSet
* [hcloud zone rrset create](hcloud_zone_rrset_create.md)	 - Create a Zone RRSet
* [hcloud zone rrset delete](hcloud_zone_rrset_delete.md)	 - Delete a Zone RRSet
* [hcloud zone rrset describe](hcloud_zone_rrset_describe.md)	 - Describe a Zone RRSet
* [hcloud zone rrset disable-protection](hcloud_zone_rrset_disable-protection.md)	 - Disable resource protection for a Zone RRSet
* [hcloud zone rrset enable-protection](hcloud_zone_rrset_enable-protection.md)	 - Enable resource protection for a Zone RRSet
* [hcloud zone rrset list](hcloud_zone_rrset_list.md)	 - List Zone RRSets
* [hcloud zone rrset remove-label](hcloud_zone_rrset_remove-label.md)	 - Remove a label from a Zone RRSet
* [hcloud zone rrset remove-records](hcloud_zone_rrset_remove-records.md)	 - Remove records from a Zone RRSet.
* [hcloud zone rrset set-records](hcloud_zone_rrset_set-records.md)	 - Set the records of a Zone RRSet

