package hcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// Action represents an action in the Hetzner Cloud.
type Action struct {
	ID           int
	Status       ActionStatus
	Command      string
	Progress     int
	Started      time.Time
	Finished     time.Time
	ErrorCode    string
	ErrorMessage string
	Resources    []*ActionResource
}

// ActionStatus represents an action's status.
type ActionStatus string

// List of action statuses.
const (
	ActionStatusRunning ActionStatus = "running"
	ActionStatusSuccess              = "success"
	ActionStatusError                = "error"
)

// ActionResource references other resources from an action.
type ActionResource struct {
	ID   int
	Type ActionResourceType
}

// ActionResourceType represents an action's resource reference type.
type ActionResourceType string

// List of action resource reference types.
const (
	ActionResourceTypeServer     ActionResourceType = "server"
	ActionResourceTypeImage                         = "image"
	ActionResourceTypeISO                           = "iso"
	ActionResourceTypeFloatingIP                    = "floating_ip"
)

func (a *Action) Error() error {
	if a.ErrorCode != "" && a.ErrorMessage != "" {
		return fmt.Errorf("%s (%s)", a.ErrorMessage, a.ErrorCode)
	}
	return nil
}

// ActionClient is a client for the actions API.
type ActionClient struct {
	client *Client
}

// GetByID retrieves an action by its ID.
func (c *ActionClient) GetByID(ctx context.Context, id int) (*Action, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/actions/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ActionGetResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		if IsError(err, ErrorCodeNotFound) {
			return nil, resp, nil
		}
		return nil, nil, err
	}
	return ActionFromSchema(body.Action), resp, nil
}

// ActionListOpts specifies options for listing actions.
type ActionListOpts struct {
	ListOpts
}

// List returns a list of actions for a specific page.
func (c *ActionClient) List(ctx context.Context, opts ActionListOpts) ([]*Action, *Response, error) {
	path := "/actions?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body schema.ActionListResponse
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	actions := make([]*Action, 0, len(body.Actions))
	for _, i := range body.Actions {
		actions = append(actions, ActionFromSchema(i))
	}
	return actions, resp, nil
}

// All returns all actions.
func (c *ActionClient) All(ctx context.Context) ([]*Action, error) {
	allActions := []*Action{}

	opts := ActionListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		actions, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allActions = append(allActions, actions...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allActions, nil
}
