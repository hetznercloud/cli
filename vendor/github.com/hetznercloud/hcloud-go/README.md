# hcloud: A Go library for the Hetzner Cloud API

[![Build status](https://travis-ci.org/hetznercloud/hcloud-go.svg?branch=master)](https://travis-ci.org/hetznercloud/hcloud-go)
[![GoDoc](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud?status.svg)](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud)

Package hcloud is a library for the Hetzner Cloud API.

The library’s documentation is available at [GoDoc](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud),
the public API documentation is available at [docs.hetzner.cloud](https://docs.hetzner.cloud/).

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/hetznercloud/hcloud-go/hcloud"
)

func main() {
    client := hcloud.NewClient(hcloud.WithToken("token"))

    server, _, err := client.Server.GetByID(context.Background(), 1)
    if err != nil {
        log.Fatalf("error retrieving server: %s\n", err)
    }
    if server != nil {
        fmt.Printf("server 1 is called %q\n", server.Name)
    } else {
        fmt.Println("server 1 not found")
    }
}
```

## License

MIT license
