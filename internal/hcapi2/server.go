package hcapi2

import (
	"context"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type ServerClient interface {
	hcloud.IServerClient
	ServerName(id int64) string
	Names() []string
	LabelKeys(idOrName string) []string
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

// ServerName obtains the name of the server with id. If the name could not
// be fetched it returns the value id converted to a string.
func (c *serverClient) ServerName(id int64) string {
	if err := c.init(); err != nil {
		return strconv.FormatInt(id, 10)
	}

	srv, ok := c.srvByID[id]
	if !ok || srv.Name == "" {
		return strconv.FormatInt(id, 10)
	}
	return srv.Name
}

// Names obtains a list of available servers. It returns nil if the
// server names could not be fetched or if there are no servers.
func (c *serverClient) Names() []string {
	if err := c.init(); err != nil || len(c.srvByID) == 0 {
		return nil
	}
	names := make([]string, len(c.srvByID))
	i := 0
	for _, srv := range c.srvByID {
		name := srv.Name
		if name == "" {
			name = strconv.FormatInt(srv.ID, 10)
		}
		names[i] = name
		i++
	}
	return names
}

// LabelKeys returns a slice containing the keys of all labels assigned
// to the Server with the passed idOrName.
func (c *serverClient) LabelKeys(idOrName string) []string {
	var srv *hcloud.Server

	if err := c.init(); err != nil || len(c.srvByID) == 0 {
		return nil
	}
	// Try to get server by ID.
	if id, err := strconv.ParseInt(idOrName, 10, 64); err != nil {
		srv = c.srvByID[id]
	}
	// If the above failed idOrName might contain a server name. If srv is not
	// nil at this point and we found something in the map, someone gave their
	// server a name containing the ID of another server.
	if v, ok := c.srvByName[idOrName]; ok && srv == nil {
		srv = v
	}
	if srv == nil || len(srv.Labels) == 0 {
		return nil
	}
	return labelKeys(srv.Labels)
}

// ServerTypeNames returns a slice of all available server types.
func (c *serverClient) ServerTypeNames() []string {
	sts, err := c.ServerTypes.All(context.Background())
	if err != nil || len(sts) == 0 {
		return nil
	}
	names := make([]string, len(sts))
	for i, st := range sts {
		names[i] = st.Name
	}
	return names
}

func (c *serverClient) init() error {
	c.once.Do(func() {
		srvs, err := c.All(context.Background())
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
