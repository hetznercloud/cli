package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/thcyron/uiprogress"
	"golang.org/x/crypto/ssh/terminal"
)

var ErrConfigPathUnknown = errors.New("config file path unknown")

type CLI struct {
	Token         string
	Endpoint      string
	Context       context.Context
	Config        *Config
	ConfigPath    string
	Debug         bool
	DebugFilePath string

	RootCommand *cobra.Command

	client *hcloud.Client

	serverNames  map[int]string
	networkNames map[int]string
}

func NewCLI() *CLI {
	cli := &CLI{
		Context:    context.Background(),
		Config:     &Config{},
		ConfigPath: DefaultConfigPath,
	}
	if s := os.Getenv("HCLOUD_CONFIG"); s != "" {
		cli.ConfigPath = s
	}
	cli.RootCommand = NewRootCommand(cli)
	return cli
}

func (c *CLI) ReadEnv() {
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
		if context := c.Config.ContextByName(s); context != nil {
			c.Config.ActiveContext = context
			c.Token = context.Token
		} else {
			log.Printf("warning: context %q specified in HCLOUD_CONTEXT does not exist\n", s)
		}
	}
}

func (c *CLI) ReadConfig() error {
	if c.ConfigPath == "" {
		return ErrConfigPathUnknown
	}

	data, err := ioutil.ReadFile(c.ConfigPath)
	if err != nil {
		return err
	}

	config, err := UnmarshalConfig(data)
	if err != nil {
		return err
	}
	if config == nil {
		return nil
	}
	c.Config = config

	if config.ActiveContext != nil {
		c.Token = config.ActiveContext.Token
	}
	if config.Endpoint != "" {
		c.Endpoint = config.Endpoint
	}

	return nil
}

func (c *CLI) WriteConfig() error {
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

func (c *CLI) wrap(f func(*CLI, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(c, cmd, args)
	}
}

func (c *CLI) Client() *hcloud.Client {
	if c.client == nil {
		opts := []hcloud.ClientOption{
			hcloud.WithToken(c.Token),
			hcloud.WithApplication("hcloud-cli", Version),
		}
		if c.Endpoint != "" {
			opts = append(opts, hcloud.WithEndpoint(c.Endpoint))
		}
		if c.Debug != false {
			if c.DebugFilePath == "" {
				opts = append(opts, hcloud.WithDebugWriter(os.Stdout))
			} else {
				writer, _ := os.Create(c.DebugFilePath)
				opts = append(opts, hcloud.WithDebugWriter(writer))
			}
		}
		pollInterval, _ := c.RootCommand.PersistentFlags().GetDuration("poll-interval")
		if pollInterval > 0 {
			opts = append(opts, hcloud.WithPollInterval(pollInterval))
		}
		c.client = hcloud.NewClient(opts...)
	}
	return c.client
}

// Terminal returns whether the CLI is run in a terminal.
func (c *CLI) Terminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func (c *CLI) ActionProgress(ctx context.Context, action *hcloud.Action) error {
	progressCh, errCh := c.Client().Action.WatchProgress(ctx, action)

	if c.Terminal() {
		progress := uiprogress.New()

		progress.Start()
		bar := progress.AddBar(100).AppendCompleted().PrependElapsed()
		bar.Empty = ' '

		for {
			select {
			case err := <-errCh:
				if err == nil {
					bar.Set(100)
				}
				progress.Stop()
				return err
			case p := <-progressCh:
				bar.Set(p)
			}
		}
	} else {
		return <-errCh
	}
}

func (c *CLI) ensureToken(cmd *cobra.Command, args []string) error {
	if c.Token == "" {
		return errors.New("no active context or token (see `hcloud context --help`)")
	}
	return nil
}

func (c *CLI) WaitForActions(ctx context.Context, actions []*hcloud.Action) error {
	for _, action := range actions {

		resources := make(map[string]int)
		for _, resource := range action.Resources {
			resources[string(resource.Type)] = resource.ID
		}

		switch action.Command {
		default:
			fmt.Printf("Waiting for action %s to have finished... ", action.Command)
		case "start_server":
			fmt.Printf("Waiting for server %d to have started... ", resources["server"])
		case "attach_volume":
			fmt.Printf("Waiting for volume %d to have been attached to server %d... ", resources["volume"], resources["server"])
		}

		_, errCh := c.Client().Action.WatchProgress(ctx, action)
		if err := <-errCh; err != nil {
			fmt.Println("failed")
			return err
		}
		fmt.Println("done")
	}

	return nil
}

func (c *CLI) GetServerName(id int) string {
	if c.serverNames == nil {
		c.serverNames = map[int]string{}
		servers, _ := c.Client().Server.All(c.Context)
		for _, server := range servers {
			c.serverNames[server.ID] = server.Name
		}
	}
	if serverName, ok := c.serverNames[id]; ok {
		return serverName
	}
	return strconv.Itoa(id)
}

func (c *CLI) GetNetworkName(id int) string {
	if c.networkNames == nil {
		c.networkNames = map[int]string{}
		networks, _ := c.Client().Network.All(c.Context)
		for _, network := range networks {
			c.networkNames[network.ID] = network.Name
		}
	}
	if networkName, ok := c.networkNames[id]; ok {
		return networkName
	}
	return strconv.Itoa(id)
}
