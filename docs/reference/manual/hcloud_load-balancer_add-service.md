## hcloud load-balancer add-service

Add a service to a Load Balancer

```
hcloud load-balancer add-service [options] (--protocol http | --protocol tcp --listen-port <1-65535> --destination-port <1-65535> | --protocol https --http-certificates <ids>) <load-balancer>
```

### Options

```
      --destination-port int                     Destination port of the service on the targets
      --health-check-http-domain string          The domain we request when performing a http health check
      --health-check-http-path string            The path we request when performing a http health check
      --health-check-http-response string        The response we expect to determine a target as healthy
      --health-check-http-status-codes strings   List of status codes we expect to determine a target as healthy
      --health-check-http-tls                    Determine if the health check should verify if the target answers with a valid TLS certificate (true, false)
      --health-check-interval duration           The interval the health check is performed (default 15s)
      --health-check-port int                    The port the health check is performed over
      --health-check-protocol string             The protocol the health check is performed over
      --health-check-retries int                 Number of retries after a health check is marked as failed (default 3)
      --health-check-timeout duration            The timeout after a health check is marked as failed (default 10s)
  -h, --help                                     help for add-service
      --http-certificates strings                IDs or names of Certificates which should be attached to this Load Balancer
      --http-cookie-lifetime duration            Sticky Sessions: Lifetime of the cookie
      --http-cookie-name string                  Sticky Sessions: Cookie Name we set
      --http-redirect-http                       Redirect all traffic on port 80 to port 443 (true, false)
      --http-sticky-sessions                     Enable Sticky Sessions (true, false)
      --listen-port int                          Listen port of the service
      --protocol string                          Protocol of the service (required)
      --proxy-protocol                           Enable proxyprotocol (true, false)
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
