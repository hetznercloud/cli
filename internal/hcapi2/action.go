package hcapi2

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ActionClient embeds the Hetzner Cloud Action ActualClient
type ActionClient interface {
	hcloud.IActionClient
}

func NewActionClient(client hcloud.IActionClient) ActionClient {
	return &actionClient{
		IActionClient: client,
	}
}

type actionClient struct {
	hcloud.IActionClient
}
