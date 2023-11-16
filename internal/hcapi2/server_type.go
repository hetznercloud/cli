package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type ServerTypeClient interface {
	hcloud.IServerTypeClient
	Names() []string
	ServerTypeName(id int64) string
	ServerTypeDescription(id int64) string
}

func NewServerTypeClient(client hcloud.IServerTypeClient) ServerTypeClient {
	return &serverTypeClient{
		IServerTypeClient: client,
	}
}

type serverTypeClient struct {
	hcloud.IServerTypeClient

	srvTypeByID map[int64]*hcloud.ServerType
	once        sync.Once
	err         error
}

// ServerTypeName obtains the name of the server type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *serverTypeClient) ServerTypeName(id int64) string {
	if err := c.init(); err != nil {
		return strconv.FormatInt(id, 10)
	}

	serverType, ok := c.srvTypeByID[id]
	if !ok || serverType.Name == "" {
		return strconv.FormatInt(id, 10)
	}
	return serverType.Name
}

// ServerTypeDescription obtains the description of the server type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *serverTypeClient) ServerTypeDescription(id int64) string {
	if err := c.init(); err != nil {
		return strconv.FormatInt(id, 10)
	}

	serverType, ok := c.srvTypeByID[id]
	if !ok || serverType.Description == "" {
		return strconv.FormatInt(id, 10)
	}
	return serverType.Description
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

func (c *serverTypeClient) init() error {
	c.once.Do(func() {
		serverTypes, err := c.All(context.Background())
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(serverTypes) == 0 {
			return
		}
		c.srvTypeByID = make(map[int64]*hcloud.ServerType, len(serverTypes))
		for _, srv := range serverTypes {
			c.srvTypeByID[srv.ID] = srv
		}
	})
	return c.err
}
