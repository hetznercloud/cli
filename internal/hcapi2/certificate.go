package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// CertificateClient embeds the Hetzner Cloud Certificate client and provides some
// additional helper functions.
type CertificateClient interface {
	hcloud.ICertificateClient
	Names() []string
	LabelKeys(string) []string
}

func NewCertificateClient(client hcloud.ICertificateClient) CertificateClient {
	return &certificateClient{
		ICertificateClient: client,
	}
}

type certificateClient struct {
	hcloud.ICertificateClient
}

// Names obtains a list of available data centers. It returns nil if
// data center names could not be fetched.
func (c *certificateClient) Names() []string {
	dcs, err := c.All(context.Background())
	if err != nil || len(dcs) == 0 {
		return nil
	}
	names := make([]string, len(dcs))
	for i, dc := range dcs {
		name := dc.Name
		if name == "" {
			name = strconv.FormatInt(dc.ID, 10)
		}
		names[i] = name
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the certificate with the passed idOrName.
func (c *certificateClient) LabelKeys(idOrName string) []string {
	certificate, _, err := c.Get(context.Background(), idOrName)
	if err != nil || certificate == nil || len(certificate.Labels) == 0 {
		return nil
	}
	return labelKeys(certificate.Labels)
}
