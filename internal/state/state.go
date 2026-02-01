package state

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

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

	if ep, err := config.OptionHetznerEndpoint.Get(c.config); err == nil && ep != "" {
		opts = append(opts, hcloud.WithHetznerEndpoint(ep))
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

		var debugWriter io.Writer
		if filePath == "" {
			debugWriter = os.Stderr
		} else {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //nolint:gosec
			if err != nil {
				return nil, err
			}
			debugWriter = f
		}

		quotedArgs := make([]string, 0, len(os.Args))
		for _, arg := range os.Args {
			quotedArgs = append(quotedArgs, fmt.Sprintf("%q", arg))
		}
		_, err = debugWriter.Write([]byte("--- Command:\n" + strings.Join(quotedArgs, " ") + "\n\n\n\n"))
		if err != nil {
			return nil, err
		}

		opts = append(opts, hcloud.WithDebugWriter(debugWriter))
	}

	pollInterval, err := config.OptionPollInterval.Get(c.config)
	if err != nil {
		return nil, err
	}

	if pollInterval > 0 {
		opts = append(opts, hcloud.WithPollOpts(hcloud.PollOpts{
			BackoffFunc: hcloud.ConstantBackoff(pollInterval),
		}))
	}

	if staticIp, err := config.OptionStaticIp.Get(c.config); err == nil && staticIp != "" {
		tr := &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			split := strings.Split(addr, ":")
			return net.Dial(network, fmt.Sprintf("%s:%s", staticIp, split[1]))
		}}
		client := &http.Client{Transport: tr}
		opts = append(opts, hcloud.WithHTTPClient(client))
	} else if err != nil {
		return nil, err
	}

	return hcapi2.NewClient(opts...), nil
}
