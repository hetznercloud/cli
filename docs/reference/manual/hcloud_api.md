## hcloud api

Make API call

```
hcloud api [options] <path>
```

### Options

```
  -d, --data string            HTTP request body content (use - to read from stdin)
  -h, --help                   help for api
  -X, --method string          HTTP method to use for API call (default "GET")
  -v, --value stringToString   key=value pairs to pass as query parameters (default [])
```

### Options inherited from parent commands

```
      --config string              Config file path (default "~/.config/hcloud/cli.toml")
      --context string             Currently active context
      --debug                      Enable debug output
      --debug-file string          File to write debug output to
      --endpoint string            Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --hetzner-endpoint string    Hetzner API endpoint (default "https://api.hetzner.com/v1")
      --http-timeout duration      Timeout for HTTP requests (default 0 = no timeout)
      --no-experimental-warnings   If true, experimental warnings are not shown
      --poll-interval duration     Interval at which to poll information, for example action progress (default 500ms)
      --quiet                      If true, only print error messages
```

### SEE ALSO

* [hcloud](hcloud.md)	 - Hetzner Cloud CLI

