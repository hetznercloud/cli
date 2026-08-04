package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// CertificateClient embeds the Hetzner Cloud Certificate client and provides some
// additional helper functions.
type CertificateClient interface {
	hcloud.ICertificateClient
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
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
func (c *certificateClient) Names(ctx context.Context) ([]string, error) {
	certificates, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(certificates, func(certificate *hcloud.Certificate) int64 { return certificate.ID }, func(certificate *hcloud.Certificate) string { return certificate.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels
// assigned to the certificate with the passed idOrName.
func (c *certificateClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	certificate, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if certificate == nil {
		return nil, nil
	}
	return labelKeys(certificate.Labels), nil
}
