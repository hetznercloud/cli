package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/thcyron/uiprogress"
	"golang.org/x/crypto/ssh/terminal"
)

var ErrConfigPathUnknown = errors.New("config file path unknown")

type CLI struct {
	Token      string
	Endpoint   string
	Context    context.Context
	Config     *Config
	ConfigPath string

	RootCommand *cobra.Command

	client *hcloud.Client
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
	if s := os.Getenv("HCLOUD_CONTEXT"); s != "" && c.Config != nil {
		c.Config.ActiveContext = c.Config.ContextByName(s)
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
	if len(actions) > 0 {
		for _, action := range actions {
			fmt.Print(actionToString(action))
			_, errCh := c.Client().Action.WatchProgress(ctx, action)
			if err := <-errCh; err != nil {
				fmt.Println("failed")
				return err
			}
			fmt.Println("done")
		}
	}
	return nil
}
