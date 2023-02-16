package state

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hetznercloud/cli/internal/hcapi"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type State struct {
	Token         string
	Endpoint      string
	Context       context.Context
	Config        *Config
	ConfigPath    string
	Debug         bool
	DebugFilePath string

	client               *hcloud.Client
	isoClient            *hcapi.ISOClient
	imageClient          *hcapi.ImageClient
	locationClient       *hcapi.LocationClient
	dataCenterClient     *hcapi.DataCenterClient
	sshKeyClient         *hcapi.SSHKeyClient
	volumeClient         *hcapi.VolumeClient
	floatingIPClient     *hcapi.FloatingIPClient
	networkClient        *hcapi.NetworkClient
	loadBalancerClient   *hcapi.LoadBalancerClient
	serverClient         *hcapi.ServerClient
	firewallClient       *hcapi.FirewallClient
	placementGroupClient *hcapi.PlacementGroupClient
}

func New() *State {
	s := &State{
		Context:    context.Background(),
		Config:     &Config{},
		ConfigPath: DefaultConfigPath,
	}
	if p := os.Getenv("HCLOUD_CONFIG"); p != "" {
		s.ConfigPath = p
	}
	return s
}

var ErrConfigPathUnknown = errors.New("config file path unknown")

func (c *State) ReadEnv() {
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
	if s := os.Getenv("HCLOUD_CONTEXT"); s != "" && c.Config != nil {
		if cfgCtx := c.Config.ContextByName(s); cfgCtx != nil {
			c.Config.ActiveContext = cfgCtx
			c.Token = cfgCtx.Token
		} else {
			log.Printf("warning: context %q specified in HCLOUD_CONTEXT does not exist\n", s)
		}
	}
}

func (c *State) ReadConfig() error {
	if c.ConfigPath == "" {
		return ErrConfigPathUnknown
	}

	data, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}

	if err = UnmarshalConfig(c.Config, data); err != nil {
		return err
	}

	if c.Config.ActiveContext != nil {
		c.Token = c.Config.ActiveContext.Token
	}
	if c.Config.Endpoint != "" {
		c.Endpoint = c.Config.Endpoint
	}

	return nil
}

func (c *State) WriteConfig() error {
	if c.ConfigPath == "" {
		return ErrConfigPathUnknown
	}
	if c.Config == nil {
		return nil
	}

	data, err := MarshalConfig(c.Config)
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
