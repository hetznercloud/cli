## hcloud context list

List contexts

### Synopsis

Displays a list of contexts.

Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=active,name (see available columns below).

Columns:
 - active
 - name

```
hcloud context list [options]
```

### Options

```
  -h, --help                 help for list
  -o, --output stringArray   output options: noheader|columns=...
```

### Options inherited from parent commands

```
      --config string            Config file path (default "~/.config/hcloud/cli.toml")
      --context string           Currently active context
      --debug                    Enable debug output
      --debug-file string        File to write debug output to
      --endpoint string          Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --poll-interval duration   Interval at which to poll information, for example action progress (default 500ms)
      --quiet                    If true, only print error messages
```

### SEE ALSO

* [hcloud context](hcloud_context.md)	 - Manage contexts
