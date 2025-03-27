## hcloud image describe

Describe an image

```
hcloud image describe [options] <image>
```

### Options

```
  -a, --architecture string   architecture of the image, default is x86 (default "x86")
  -h, --help                  help for describe
  -o, --output stringArray    output options: json|yaml|format
```

### Options inherited from parent commands

```
      --config string            Config file path (default "/Users/paul/.config/hcloud/cli.toml")
      --context string           Currently active context
      --debug                    Enable debug output
      --debug-file string        File to write debug output to
      --endpoint string          Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --poll-interval duration   Interval at which to poll information, for example action progress (default 500ms)
      --quiet                    If true, only print error messages
```

### SEE ALSO

* [hcloud image](hcloud_image.md)	 - Manage images
