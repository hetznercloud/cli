## hcloud zone rrset

[experimental] Manage Zone RRSets (records)

### Synopsis

For more details, see the documentation for Zones https://docs.hetzner.cloud/reference/cloud#zones or Zone RRSets https://docs.hetzner.cloud/reference/cloud#zone-rrsets.

Experimental: DNS API is in beta, breaking changes may occur within minor releases.
See https://docs.hetzner.cloud/changelog#2025-10-07-dns-beta for more details.


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
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud zone](hcloud_zone.md)	 - [experimental] Manage DNS Zones and Zone RRSets (records)
* [hcloud zone rrset add-label](hcloud_zone_rrset_add-label.md)	 - [experimental] Add a label to a Zone RRSet
* [hcloud zone rrset add-records](hcloud_zone_rrset_add-records.md)	 - [experimental] Add records to a Zone RRSet
* [hcloud zone rrset change-ttl](hcloud_zone_rrset_change-ttl.md)	 - [experimental] Changes the Time To Live (TTL) of a Zone RRSet
* [hcloud zone rrset create](hcloud_zone_rrset_create.md)	 - [experimental] Create a Zone RRSet
* [hcloud zone rrset delete](hcloud_zone_rrset_delete.md)	 - [experimental] Delete a Zone RRSet
* [hcloud zone rrset describe](hcloud_zone_rrset_describe.md)	 - [experimental] Describe a Zone RRSet
* [hcloud zone rrset disable-protection](hcloud_zone_rrset_disable-protection.md)	 - [experimental] Disable resource protection for a Zone RRSet
* [hcloud zone rrset enable-protection](hcloud_zone_rrset_enable-protection.md)	 - [experimental] Enable resource protection for a Zone RRSet
* [hcloud zone rrset list](hcloud_zone_rrset_list.md)	 - [experimental] List Zone RRSets
* [hcloud zone rrset remove-label](hcloud_zone_rrset_remove-label.md)	 - [experimental] Remove a label from a Zone RRSet
* [hcloud zone rrset remove-records](hcloud_zone_rrset_remove-records.md)	 - [experimental] Remove records from a Zone RRSet.
* [hcloud zone rrset set-records](hcloud_zone_rrset_set-records.md)	 - [experimental] Set the records of a Zone RRSet
