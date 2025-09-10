## hcloud storage-box enable-snapshot-plan

Enable automatic snapshots for a Storage Box

```
hcloud storage-box enable-snapshot-plan [options] <storage-box>
```

### Options

```
      --day-of-month int     Day of the month the Snapshot Plan should be executed on. Not specified means every day
      --day-of-week string   Day of the week the Snapshot Plan should be executed on. Not specified means every day
  -h, --help                 help for enable-snapshot-plan
      --hour int             Hour the Snapshot Plan should be executed on (UTC). Not specified means every hour
      --max-snapshots int    Maximum amount of Snapshots that should be created by this Snapshot Plan
      --minute int           Minute the Snapshot Plan should be executed on (UTC). Not specified means every minute
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

* [hcloud storage-box](hcloud_storage-box.md)	 - Manage Storage Boxes
