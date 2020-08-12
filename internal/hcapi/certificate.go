package hcapi

import (
	"context"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// CertificateClient embeds the Hetzner Cloud Certificate client and provides
// some additional helper functions.
type CertificateClient struct {
	*hcloud.CertificateClient
}

// CertificateNames obtains a list of available certificates. It returns nil if
// the certificate names could not be fetched or there were no certificates.
func (c *CertificateClient) CertificateNames() []string {
	certs, err := c.All(context.Background())
	if err != nil || len(certs) == 0 {
		return nil
	}
	names := make([]string, len(certs))
	for i, cert := range certs {
		name := cert.Name
		if name == "" {
			name = strconv.Itoa(cert.ID)
		}
		names[i] = name
	}
	return names
}

// CertificateLabelKeys returns a slice containing the keys of all labels
// assigned to the certificate with the passed idOrName.
func (c *CertificateClient) CertificateLabelKeys(idOrName string) []string {
	cert, _, err := c.Get(context.Background(), idOrName)
	if err != nil || cert == nil || len(cert.Labels) == 0 {
		return nil
	}
	return lkeys(cert.Labels)
}
