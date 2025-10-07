## hcloud zone

[experimental] Manage DNS Zones and Zone RRSets (records)

### Synopsis

For more details, see the documentation for Zones https://docs.hetzner.cloud/reference/cloud#zones or Zone RRSets https://docs.hetzner.cloud/reference/cloud#zone-rrsets.

Experimental: DNS API is in beta, breaking changes may occur within minor releases.
See https://docs.hetzner.cloud/changelog#2025-10-07-dns-beta for more details.


### Options

```
  -h, --help   help for zone
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

* [hcloud](hcloud.md)	 - Hetzner Cloud CLI
* [hcloud zone add-label](hcloud_zone_add-label.md)	 - [experimental] Add a label to a Zone
* [hcloud zone add-records](hcloud_zone_add-records.md)	 - [experimental] Add records to a Zone RRSet
* [hcloud zone change-primary-nameservers](hcloud_zone_change-primary-nameservers.md)	 - [experimental] Changes the primary nameservers of a secondary Zone
* [hcloud zone change-ttl](hcloud_zone_change-ttl.md)	 - [experimental] Changes the default Time To Live (TTL) of a Zone
* [hcloud zone create](hcloud_zone_create.md)	 - [experimental] Create a Zone
* [hcloud zone delete](hcloud_zone_delete.md)	 - [experimental] Delete a Zone
* [hcloud zone describe](hcloud_zone_describe.md)	 - [experimental] Describe a Zone
* [hcloud zone disable-protection](hcloud_zone_disable-protection.md)	 - [experimental] Disable resource protection for a Zone
* [hcloud zone enable-protection](hcloud_zone_enable-protection.md)	 - [experimental] Enable resource protection for a Zone
* [hcloud zone export-zonefile](hcloud_zone_export-zonefile.md)	 - [experimental] Returns a generated Zone file in BIND (RFC 1034/1035) format
* [hcloud zone import-zonefile](hcloud_zone_import-zonefile.md)	 - [experimental] Imports a zone file, replacing all Zone RRSets
* [hcloud zone list](hcloud_zone_list.md)	 - [experimental] List Zones
* [hcloud zone remove-label](hcloud_zone_remove-label.md)	 - [experimental] Remove a label from a Zone
* [hcloud zone remove-records](hcloud_zone_remove-records.md)	 - [experimental] Remove records from a Zone RRSet.
* [hcloud zone rrset](hcloud_zone_rrset.md)	 - [experimental] Manage Zone RRSets (records)
* [hcloud zone set-records](hcloud_zone_set-records.md)	 - [experimental] Set the records of a Zone RRSet
