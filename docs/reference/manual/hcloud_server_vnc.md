## hcloud server vnc

[experimental] Forward the Server VNC console to a local port

### Synopsis

Forward the Server VNC console to a local port, so it can be accessed with any VNC client.

Note that the local VNC server has no password protection.

Experimental: The server vnc command is experimental, breaking changes may occur within minor releases.


```
hcloud server vnc [options] <server>
```

### Options

```
      --address string   Local address to listen on (default "127.0.0.1")
  -h, --help             help for vnc
  -p, --port int         Local port to listen on (0 picks a random free port) (default 5900)
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

* [hcloud server](hcloud_server.md)	 - Manage Servers

