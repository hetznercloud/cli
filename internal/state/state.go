package state

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/version"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type State interface {
	context.Context

	TokenEnsurer
	ActionWaiter

	Client() hcapi2.Client
	Config() *Config
}

type state struct {
	context.Context

	token         string
	endpoint      string
	debug         bool
	debugFilePath string
	client        hcapi2.Client
	hcloudClient  *hcloud.Client
	config        *Config
}

func New(cfg *Config) (State, error) {
	var (
		token    string
		endpoint string
	)
	if cfg.ActiveContext != nil {
		token = cfg.ActiveContext.Token
	}
	if cfg.Endpoint != "" {
		endpoint = cfg.Endpoint
	}

	s := &state{
		Context:  context.Background(),
		config:   cfg,
		token:    token,
		endpoint: endpoint,
	}

	s.readEnv()
	s.hcloudClient = s.newClient()
	s.client = hcapi2.NewClient(s.hcloudClient)
	return s, nil
}

func ReadConfig(path string) (*Config, error) {
	cfg := &Config{Path: path}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = UnmarshalConfig(cfg, data); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *state) Client() hcapi2.Client {
	return c.client
}

func (c *state) Config() *Config {
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
			c.config.ActiveContext = cfgCtx
			c.token = cfgCtx.Token
		} else {
			log.Printf("warning: context %q specified in HCLOUD_CONTEXT does not exist\n", s)
		}
	}
}

func (c *state) newClient() *hcloud.Client {
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
	// TODO Somehow pass here
	// pollInterval, _ := c.RootCommand.PersistentFlags().GetDuration("poll-interval")
	pollInterval := 500 * time.Millisecond
	if pollInterval > 0 {
		opts = append(opts, hcloud.WithPollInterval(pollInterval))
	}
	return hcloud.NewClient(opts...)
}
