package state

import (
	"context"
	"os"

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
	opts := []hcloud.ClientOption{
		hcloud.WithToken(config.OptionToken.Get(c.config)),
		hcloud.WithApplication("hcloud-cli", version.Version),
	}

	if ep := config.OptionEndpoint.Get(c.config); ep != "" {
		opts = append(opts, hcloud.WithEndpoint(ep))
	}
	if config.OptionDebug.Get(c.config) {
		if filePath := config.OptionDebugFile.Get(c.config); filePath == "" {
			opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
		} else {
			writer, _ := os.Create(filePath)
			opts = append(opts, hcloud.WithDebugWriter(writer))
		}
	}
	pollInterval := config.OptionPollInterval.Get(c.config)
	if pollInterval > 0 {
		opts = append(opts, hcloud.WithBackoffFunc(hcloud.ConstantBackoff(pollInterval)))
	}

	return hcapi2.NewClient(opts...)
}
