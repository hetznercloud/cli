package hcapi2

import "context"

type ServerTypeClient interface {
	ServerTypeClientBase
	Names() []string
}

func NewServerTypeClient(client ServerTypeClientBase) ServerTypeClient {
	return &serverTypeClient{
		ServerTypeClientBase: client,
	}
}

type serverTypeClient struct {
	ServerTypeClientBase
}

// Names returns a slice of all available server types.
func (c *serverTypeClient) Names() []string {
	sts, err := c.All(context.Background())
	if err != nil || len(sts) == 0 {
		return nil
	}
	names := make([]string, len(sts))
	for i, st := range sts {
		names[i] = st.Name
	}
	return names
}
