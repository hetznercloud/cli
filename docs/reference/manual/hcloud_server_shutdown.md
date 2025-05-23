## hcloud server shutdown

Shutdown a server

### Synopsis

Shuts down a Server gracefully by sending an ACPI shutdown request. The Server operating system must support ACPI and react to the request, otherwise the Server will not shut down. Use the --wait flag to wait for the Server to shut down before returning.

```
hcloud server shutdown [options] <server>
```

### Options

```
  -h, --help                    help for shutdown
      --wait                    Wait for the Server to shut down before exiting (true, false)
      --wait-timeout duration   Timeout for waiting for off state after shutdown (default 30s)
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

* [hcloud server](hcloud_server.md)	 - Manage Servers
