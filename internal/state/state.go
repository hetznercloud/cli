package state

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil/terminal"
	"github.com/hetznercloud/cli/internal/version"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type State interface {
	TokenEnsurer
	ActionWaiter
	io.Closer

	Client() hcapi2.Client
	Config() config.Config
	Terminal() terminal.Terminal
	Logger() *slog.Logger
}

type state struct {
	client  hcapi2.Client
	config  config.Config
	term    terminal.Terminal
	stderr  io.Writer
	logger  *slog.Logger
	closers []io.Closer
}

type Options struct {
	Stderr   io.Writer
	Terminal terminal.Terminal
}

func New(cfg config.Config, opts Options) (State, error) {
	if opts.Stderr == nil {
		opts.Stderr = io.Discard
	}
	if opts.Terminal == nil {
		opts.Terminal = terminal.DefaultTerminal{}
	}

	s := &state{
		config: cfg,
		term:   opts.Terminal,
		stderr: opts.Stderr,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	var err error
	s.client, err = s.newClient()
	if err != nil {
		return nil, errors.Join(err, s.Close())
	}
	return s, nil
}

func (c *state) Close() error {
	errs := make([]error, 0, len(c.closers))
	for _, closer := range c.closers {
		errs = append(errs, closer.Close())
	}
	c.closers = nil
	return errors.Join(errs...)
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

func (c *state) Logger() *slog.Logger {
	return c.logger
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
			debugWriter = c.stderr
		} else {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				return nil, err
			}
			debugWriter = f
			c.closers = append(c.closers, f)
		}

		_, err = fmt.Fprintf(debugWriter, "--- hcloud debug session (%s) ---\n\n", time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return nil, err
		}
		c.logger = slog.New(slog.NewJSONHandler(debugWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))

		opts = append(opts, hcloud.WithDebugWriter(debugWriter))
	}

	pollInterval, err := config.OptionPollInterval.Get(c.config)
	if err != nil {
		return nil, err
	}

	if pollInterval > 0 && config.OptionPollInterval.Changed(c.config) {
		opts = append(opts, hcloud.WithPollOpts(hcloud.PollOpts{
			BackoffFunc: hcloud.ConstantBackoff(pollInterval),
		}))
	} else {
		opts = append(opts, hcloud.WithPollOpts(hcloud.PollOpts{
			BackoffFunc: customPollBackoffFunc(),
		}))
	}

	httpTimeout, err := config.OptionHTTPTimeout.Get(c.config)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	client.Timeout = httpTimeout
	opts = append(opts, hcloud.WithHTTPClient(client))

	return hcapi2.NewClient(opts...), nil
}

func customPollBackoffFunc() hcloud.BackoffFunc {
	const (
		initialPollInterval = 500 * time.Millisecond
		backoffMultiplier   = 1.5
		maximumPollInterval = 2500 * time.Millisecond
		constantPollRetries = 10
	)

	constantFunc := hcloud.ConstantBackoff(initialPollInterval)

	exponentialFunc := hcloud.ExponentialBackoffWithOpts(hcloud.ExponentialBackoffOpts{
		Base:       initialPollInterval,
		Multiplier: backoffMultiplier,
		Cap:        maximumPollInterval,
	})

	// Poll every 500ms for the first 5s, then use an exponential backoff capped to 2500ms
	return func(retries int) time.Duration {
		if retries < constantPollRetries {
			return constantFunc(retries)
		}
		return exponentialFunc(retries - constantPollRetries)
	}
}
