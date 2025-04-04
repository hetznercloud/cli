## hcloud certificate create

Create or upload a Certificate

```
hcloud certificate create [options] --name <name> (--type managed --domain <domain> | --type uploaded --cert-file <file> --key-file <file>)
```

### Options

```
      --cert-file string       File containing the PEM encoded certificate (required if type is uploaded)
      --domain strings         One or more domains the Certificate is valid for.
  -h, --help                   help for create
      --key-file string        File containing the PEM encoded private key for the certificate (required if type is uploaded)
      --label stringToString   User-defined labels ('key=value') (can be specified multiple times) (default [])
      --name string            Certificate name (required)
  -o, --output stringArray     output options: json|yaml
  -t, --type string            Type of Certificate to create. Valid choices: uploaded, managed (default "uploaded")
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

* [hcloud certificate](hcloud_certificate.md)	 - Manage Certificates
