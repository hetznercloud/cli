package hcapi2

// NewRDNSClient embeds the Hetzner Cloud rdns client.
type RDNSClient interface {
	RDNSClientBase
}

func NewRDNSClient(client RDNSClientBase) RDNSClient {
	return &rdnsClient{
		RDNSClientBase: client,
	}
}

type rdnsClient struct {
	RDNSClientBase
}
