package state

import (
	"context"
	"os"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil/terminal"
	"github.com/hetznercloud/cli/internal/version"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type State interface {
	context.Context

	TokenEnsurer
	ActionWaiter

	Client() hcapi2.Client
	Config() config.Config
	Terminal() terminal.Terminal
}

type state struct {
	context.Context

	client hcapi2.Client
	config config.Config
	term   terminal.Terminal
}

func New(cfg config.Config) (State, error) {
	s := &state{
		Context: context.Background(),
		config:  cfg,
		term:    terminal.DefaultTerminal{},
	}

	var err error
	s.client, err = s.newClient()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *state) Client() hcapi2.Client {
	return c.client
}

func (c *state) Config() config.Config {
	return c.config
}

func (c *state) Terminal() terminal.Terminal {
	return c.term
}

func (c *state) newClient() (hcapi2.Client, error) {
	tok, err := config.OptionToken.Get(c.config)
	if err != nil {
		return nil, err
	}

	opts := []hcloud.ClientOption{
		hcloud.WithToken(tok),
		hcloud.WithApplication("hcloud-cli", version.Version),
	}

	if ep, err := config.OptionEndpoint.Get(c.config); err == nil && ep != "" {
		opts = append(opts, hcloud.WithEndpoint(ep))
	} else if err != nil {
		return nil, err
	}

	debug, err := config.OptionDebug.Get(c.config)
	if err != nil {
		return nil, err
	}

	if debug {
		filePath, err := config.OptionDebugFile.Get(c.config)
		if err != nil {
			return nil, err
		}

		if filePath == "" {
			opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
		} else {
			writer, _ := os.Create(filePath)
			opts = append(opts, hcloud.WithDebugWriter(writer))
		}
	}

	pollInterval, err := config.OptionPollInterval.Get(c.config)
	if err != nil {
		return nil, err
	}

	if pollInterval > 0 {
		opts = append(opts, hcloud.WithBackoffFunc(hcloud.ConstantBackoff(pollInterval)))
	}

	return hcapi2.NewClient(opts...), nil
}
