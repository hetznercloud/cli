package state

import (
	"context"
	"log"
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

	token         string
	endpoint      string
	debug         bool
	debugFilePath string
	client        hcapi2.Client
	config        config.Config
}

func New(cfg config.Config) (State, error) {
	var (
		token    string
		endpoint string
	)
	if ctx := cfg.ActiveContext(); ctx != nil {
		token = ctx.Token
	}
	if ep := cfg.Endpoint(); ep != "" {
		endpoint = ep
	}

	s := &state{
		Context:  context.Background(),
		config:   cfg,
		token:    token,
		endpoint: endpoint,
	}

	s.readEnv()
	s.client = s.newClient()
	return s, nil
}

func (c *state) Client() hcapi2.Client {
	return c.client
}

func (c *state) Config() config.Config {
	return c.config
}

func (c *state) readEnv() {
	if s := os.Getenv("HCLOUD_TOKEN"); s != "" {
		c.token = s
	}
	if s := os.Getenv("HCLOUD_ENDPOINT"); s != "" {
		c.endpoint = s
	}
	if s := os.Getenv("HCLOUD_DEBUG"); s != "" {
		c.debug = true
	}
	if s := os.Getenv("HCLOUD_DEBUG_FILE"); s != "" {
		c.debugFilePath = s
	}
	if s := os.Getenv("HCLOUD_CONTEXT"); s != "" && c.config != nil {
		if cfgCtx := c.config.ContextByName(s); cfgCtx != nil {
			c.config.SetActiveContext(cfgCtx)
			c.token = cfgCtx.Token
		} else {
			log.Printf("warning: context %q specified in HCLOUD_CONTEXT does not exist\n", s)
		}
	}
}

func (c *state) newClient() hcapi2.Client {
	opts := []hcloud.ClientOption{
		hcloud.WithToken(c.token),
		hcloud.WithApplication("hcloud-cli", version.Version),
	}
	if c.endpoint != "" {
		opts = append(opts, hcloud.WithEndpoint(c.endpoint))
	}
	if c.debug {
		if c.debugFilePath == "" {
			opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
		} else {
			writer, _ := os.Create(c.debugFilePath)
			opts = append(opts, hcloud.WithDebugWriter(writer))
		}
	}
	return hcapi2.NewClient(opts...)
}
