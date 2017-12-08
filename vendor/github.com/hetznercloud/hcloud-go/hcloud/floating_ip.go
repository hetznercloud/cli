package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// FloatingIP represents a Floating IP in the Hetzner Cloud.
type FloatingIP struct {
	ID           int
	Description  string
	IP           string
	Type         FloatingIPType
	Server       *Server
	DNSPtr       map[string]string
	HomeLocation *Location
}

// FloatingIPType represents the type of a Floating IP.
type FloatingIPType string

// Floating IP types.
const (
	FloatingIPTypeIPv4 FloatingIPType = "ipv4"
	FloatingIPTypeIPv6                = "ipv6"
)

// FloatingIPClient is a client for the Floating IP API.
type FloatingIPClient struct {
	client *Client
}

// Get retrieves a Floating IP.
func (c *FloatingIPClient) Get(ctx context.Context, id int) (*FloatingIP, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/floating_ips/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.FloatingIPGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}
	return FloatingIPFromSchema(body.FloatingIP), resp, nil
}

// FloatingIPListOpts specifies options for listing Floating IPs.
type FloatingIPListOpts struct {
	ListOpts
}

// List returns a list of Floating IPs for a specific page.
func (c *FloatingIPClient) List(ctx context.Context, opts FloatingIPListOpts) ([]*FloatingIP, *Response, error) {
	path := "/floating_ips?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.FloatingIPListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	floatingIPs := make([]*FloatingIP, 0, len(body.FloatingIPs))
	for _, s := range body.FloatingIPs {
		floatingIPs = append(floatingIPs, FloatingIPFromSchema(s))
	}
	return floatingIPs, resp, nil
}

// All returns all Floating IPs.
func (c *FloatingIPClient) All(ctx context.Context) ([]*FloatingIP, error) {
	allFloatingIPs := []*FloatingIP{}

	opts := FloatingIPListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		floatingIPs, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allFloatingIPs = append(allFloatingIPs, floatingIPs...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allFloatingIPs, nil
}

// FloatingIPCreateOpts specifies options for creating a Floating IP.
type FloatingIPCreateOpts struct {
	Type         FloatingIPType
	HomeLocation *Location
	Server       *Server
	Description  *string
}

// Validate checks if options are valid.
func (o FloatingIPCreateOpts) Validate() error {
	switch o.Type {
	case FloatingIPTypeIPv4, FloatingIPTypeIPv6:
		break
	default:
		return errors.New("missing or invalid type")
	}
	if o.HomeLocation == nil && o.Server == nil {
		return errors.New("one of home location or server is required")
	}
	return nil
}

// FloatingIPCreateResult is the result of creating a Floating IP.
type FloatingIPCreateResult struct {
	FloatingIP *FloatingIP
	Action     *Action
}

// Create creates a Floating IP.
func (c *FloatingIPClient) Create(ctx context.Context, opts FloatingIPCreateOpts) (FloatingIPCreateResult, *Response, error) {
	if err := opts.Validate(); err != nil {
		return FloatingIPCreateResult{}, nil, err
	}

	reqBody := schema.FloatingIPCreateRequest{
		Type:        string(opts.Type),
		Description: opts.Description,
	}
	if opts.HomeLocation != nil {
		reqBody.HomeLocation = String(opts.HomeLocation.Name)
	}
	if opts.Server != nil {
		reqBody.Server = Int(opts.Server.ID)
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return FloatingIPCreateResult{}, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/floating_ips", bytes.NewReader(reqBodyData))
	if err != nil {
		return FloatingIPCreateResult{}, nil, err
	}

	var respBody schema.FloatingIPCreateResponse
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return FloatingIPCreateResult{}, resp, err
	}
	var action *Action
	if respBody.Action != nil {
		action = ActionFromSchema(*respBody.Action)
	}
	return FloatingIPCreateResult{
		FloatingIP: FloatingIPFromSchema(respBody.FloatingIP),
		Action:     action,
	}, resp, nil
}
