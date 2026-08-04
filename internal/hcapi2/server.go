package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type ServerClient interface {
	hcloud.IServerClient
	ServerName(context.Context, int64) (string, error)
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
}

func NewServerClient(client *hcloud.ServerClient) ServerClient {
	return &serverClient{
		IServerClient: client,
	}
}

// ServerClient embeds the Hetzner Cloud Server client and provides some
// additional helper functions.
type serverClient struct {
	hcloud.IServerClient

	ServerTypes *hcloud.ServerTypeClient

	srvByID   map[int64]*hcloud.Server
	srvByName map[string]*hcloud.Server

	once sync.Once
	err  error
}

// ServerName obtains the name of the server with id. It returns the numeric ID
// when the API response does not contain a matching named server.
func (c *serverClient) ServerName(ctx context.Context, id int64) (string, error) {
	if err := c.init(ctx); err != nil {
		return "", err
	}

	srv, ok := c.srvByID[id]
	if !ok || srv.Name == "" {
		return strconv.FormatInt(id, 10), nil
	}
	return srv.Name, nil
}

// Names obtains a list of available servers. It returns nil if the
// server names could not be fetched or if there are no servers.
func (c *serverClient) Names(ctx context.Context) ([]string, error) {
	servers, err := c.All(ctx)
	if err != nil {
		return nil, err
	}
	return resourceNames(servers, func(server *hcloud.Server) int64 { return server.ID }, func(server *hcloud.Server) string { return server.Name }), nil
}

// LabelKeys returns a slice containing the keys of all labels assigned
// to the Server with the passed idOrName.
func (c *serverClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	server, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, nil
	}
	return labelKeys(server.Labels), nil
}

func (c *serverClient) init(ctx context.Context) error {
	c.once.Do(func() {
		srvs, err := c.All(ctx)
		if err != nil {
			c.err = err
		}
		if c.err != nil || len(srvs) == 0 {
			return
		}
		c.srvByID = make(map[int64]*hcloud.Server, len(srvs))
		c.srvByName = make(map[string]*hcloud.Server, len(srvs))
		for _, srv := range srvs {
			c.srvByID[srv.ID] = srv
			c.srvByName[srv.Name] = srv
		}
	})
	return c.err
}
