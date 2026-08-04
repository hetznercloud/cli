package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type ServerTypeClient interface {
	hcloud.IServerTypeClient
	Names(context.Context) ([]string, error)
	ServerTypeName(context.Context, int64) (string, error)
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

// ServerTypeName obtains the name of the server type with id. It returns the
// numeric ID when the API response does not contain a matching named type.
func (c *serverTypeClient) ServerTypeName(ctx context.Context, id int64) (string, error) {
	if err := c.init(ctx); err != nil {
		return "", err
	}

	serverType, ok := c.srvTypeByID[id]
	if !ok || serverType.Name == "" {
		return strconv.FormatInt(id, 10), nil
	}
	return serverType.Name, nil
}

// Names returns a slice of all available server types.
func (c *serverTypeClient) Names(ctx context.Context) ([]string, error) {
	serverTypes, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(serverTypes))
	for i, st := range serverTypes {
		names[i] = st.Name
	}
	return names, nil
}

func (c *serverTypeClient) init(ctx context.Context) error {
	c.once.Do(func() {
		serverTypes, err := c.All(ctx)
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
