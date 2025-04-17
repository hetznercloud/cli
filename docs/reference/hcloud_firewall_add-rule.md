## hcloud firewall add-rule

Add a single rule to a firewall

```
hcloud firewall add-rule [options] (--direction in --source-ips <ips> | --direction out --destination-ips <ips>) (--protocol <tcp|udp> --port <port> | --protocol <icmp|esp|gre>) <firewall>
```

### Options

```
      --description string        Description of the Firewall rule
      --destination-ips strings   Destination IPs (CIDR Notation) (required when direction is out)
      --direction string          Direction (in, out) (required)
  -h, --help                      help for add-rule
      --port string               Port to which traffic will be allowed, only applicable for protocols TCP and UDP, you can specify port ranges, sample: 80-85
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
