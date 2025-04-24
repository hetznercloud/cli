## hcloud firewall delete-rule

Delete a single rule from a Firewall

```
hcloud firewall delete-rule [options] (--direction in --source-ips <ips> | --direction out --destination-ips <ips>) (--protocol <tcp|udp> --port <port> | --protocol <icmp|esp|gre>) <firewall>
```

### Options

```
      --description string        Description of the Firewall rule
      --destination-ips strings   Destination IPs (CIDR Notation) (required when direction is out)
      --direction string          Direction (in, out) (required)
  -h, --help                      help for delete-rule
      --port string               Port to which traffic will be allowed, only applicable for protocols TCP and UDP
      --protocol string           Protocol (icmp, esp, gre, udp or tcp) (required)
      --source-ips strings        Source IPs (CIDR Notation) (required when direction is in)
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

* [hcloud firewall](hcloud_firewall.md)	 - Manage Firewalls
