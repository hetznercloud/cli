package state

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/version"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type State interface {
	TokenEnsurer
	ActionWaiter
	context.Context
	hcapi2.Client

	WriteConfig() error
	Config() *Config
}

type state struct {
	context.Context
	hcapi2.Client

	Token         string
	Endpoint      string
	ConfigPath    string
	Debug         bool
	DebugFilePath string

	hcloudClient *hcloud.Client
	config       *Config
}

func New() (State, error) {
	s := &state{
		Context:    context.Background(),
		config:     &Config{},
		ConfigPath: DefaultConfigPath,
	}
	if p := os.Getenv("HCLOUD_CONFIG"); p != "" {
		s.ConfigPath = p
	}
	if err := s.readConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config file %q: %s\n", s.ConfigPath, err)
	}
	s.readEnv()
	s.hcloudClient = s.newClient()
	s.Client = hcapi2.NewClient(s.hcloudClient)
	return s, nil
}

func (c *state) Config() *Config {
	return c.config
}

func (c *state) readEnv() {
	if s := os.Getenv("HCLOUD_TOKEN"); s != "" {
		c.Token = s
	}
	if s := os.Getenv("HCLOUD_ENDPOINT"); s != "" {
		c.Endpoint = s
	}
	if s := os.Getenv("HCLOUD_DEBUG"); s != "" {
		c.Debug = true
	}
	if s := os.Getenv("HCLOUD_DEBUG_FILE"); s != "" {
		c.DebugFilePath = s
	}
	if s := os.Getenv("HCLOUD_CONTEXT"); s != "" && c.config != nil {
		if cfgCtx := c.config.ContextByName(s); cfgCtx != nil {
			c.config.ActiveContext = cfgCtx
			c.Token = cfgCtx.Token
		} else {
			log.Printf("warning: context %q specified in HCLOUD_CONTEXT does not exist\n", s)
		}
	}
}

func (c *state) readConfig() error {
	_, err := os.Stat(c.ConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}
	if err = UnmarshalConfig(c.config, data); err != nil {
		return err
	}

	if c.config.ActiveContext != nil {
		c.Token = c.config.ActiveContext.Token
	}
	if c.config.Endpoint != "" {
		c.Endpoint = c.config.Endpoint
	}
	return nil
}

func (c *state) newClient() *hcloud.Client {
	opts := []hcloud.ClientOption{
		hcloud.WithToken(c.Token),
		hcloud.WithApplication("hcloud-cli", version.Version),
	}
	if c.Endpoint != "" {
		opts = append(opts, hcloud.WithEndpoint(c.Endpoint))
	}
	if c.Debug {
		if c.DebugFilePath == "" {
			opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
		} else {
			writer, _ := os.Create(c.DebugFilePath)
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

func (c *state) WriteConfig() error {
	if c.config == nil {
		return nil
	}

	data, err := MarshalConfig(c.config)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(c.ConfigPath), 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.ConfigPath, data, 0600); err != nil {
		return err
	}
	return nil
}
