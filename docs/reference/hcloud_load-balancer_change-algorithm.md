## hcloud load-balancer change-algorithm

Changes the algorithm of a Load Balancer

```
hcloud load-balancer change-algorithm --algorithm-type <round_robin|least_connections> <load-balancer>
```

### Options

```
      --algorithm-type string   New Load Balancer algorithm (round_robin, least_connections) (required)
  -h, --help                    help for change-algorithm
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

* [hcloud load-balancer](hcloud_load-balancer.md)	 - Manage Load Balancers
