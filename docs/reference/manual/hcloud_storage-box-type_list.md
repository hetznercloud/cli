## hcloud storage-box-type list

[experimental] List Storage Box Types

### Synopsis

Displays a list of Storage Box Types.

Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=automatic_snapshot_limit,deprecated (see available columns below).

Columns:
 - automatic_snapshot_limit
 - deprecated
 - description
 - id
 - name
 - size
 - snapshot_limit
 - subaccounts_limit

Experimental: Storage Box support is experimental, breaking changes may occur within minor releases.
See https://github.com/hetznercloud/cli/issues/1202 for more details.


```
hcloud storage-box-type list [options]
```

### Options

```
  -h, --help                 help for list
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
      --hetzner-endpoint string    Hetzner API endpoint (default "https://api.hetzner.com/v1")
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud storage-box-type](hcloud_storage-box-type.md)	 - [experimental] View Storage Box Types
