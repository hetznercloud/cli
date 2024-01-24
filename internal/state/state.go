package state

import (
	"context"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/version"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type State interface {
	context.Context

	TokenEnsurer
	ActionWaiter

	Client() hcapi2.Client
	Config() config.Config
}

type state struct {
	context.Context

	client hcapi2.Client
	config config.Config
}

func New(cfg config.Config) (State, error) {
	s := &state{
		Context: context.Background(),
		config:  cfg,
	}

	s.client = s.newClient()
	return s, nil
}

func (c *state) Client() hcapi2.Client {
	return c.client
}

func (c *state) Config() config.Config {
	return c.config
}

func (c *state) newClient() hcapi2.Client {
	opts := config.GetHcloudOpts(c.Config())
	opts = append(opts, hcloud.WithApplication("hcloud-cli", version.Version))
	return hcapi2.NewClient(opts...)
}
