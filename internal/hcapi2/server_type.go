package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type ServerTypeClient interface {
	ServerTypeClientBase
	Names() []string
	ServerTypeName(id int) string
	ServerTypeDescription(id int) string
}

func NewServerTypeClient(client ServerTypeClientBase) ServerTypeClient {
	return &serverTypeClient{
		ServerTypeClientBase: client,
	}
}

type serverTypeClient struct {
	ServerTypeClientBase

	srvTypeByID map[int]*hcloud.ServerType
	once        sync.Once
	err         error
}

// ServerTypeName obtains the name of the server type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *serverTypeClient) ServerTypeName(id int) string {
	if err := c.init(); err != nil {
		return strconv.Itoa(id)
	}

	serverType, ok := c.srvTypeByID[id]
	if !ok || serverType.Name == "" {
		return strconv.Itoa(id)
	}
	return serverType.Name
}

// ServerTypeDescription obtains the description of the server type with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *serverTypeClient) ServerTypeDescription(id int) string {
	if err := c.init(); err != nil {
		return strconv.Itoa(id)
	}

	serverType, ok := c.srvTypeByID[id]
	if !ok || serverType.Description == "" {
		return strconv.Itoa(id)
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
		c.srvTypeByID = make(map[int]*hcloud.ServerType, len(serverTypes))
		for _, srv := range serverTypes {
			c.srvTypeByID[srv.ID] = srv
		}
	})
	return c.err
}
