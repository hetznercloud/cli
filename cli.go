package cli

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/thcyron/uiprogress"
	"golang.org/x/crypto/ssh/terminal"
)

type CLI struct {
	Token    string
	Endpoint string
	Context  context.Context
	Config   *Config

	RootCommand *cobra.Command

	client *hcloud.Client
}

func NewCLI() *CLI {
	cli := &CLI{
		Context: context.Background(),
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
}

func (c *CLI) ReadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config, err := UnmarshalConfig(data)
	if err != nil {
		return err
	}
	c.Config = config

	if config.Token != "" {
		c.Token = config.Token
	}
	if config.Endpoint != "" {
		c.Endpoint = config.Endpoint
	}

	return nil
}

func (c *CLI) WriteConfig(path string) error {
	if c.Config == nil {
		return nil
	}

	data, err := MarshalConfig(c.Config)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
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
		}
		if c.Endpoint != "" {
			opts = append(opts, hcloud.WithEndpoint(c.Endpoint))
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
	errCh, progressCh := waitAction(ctx, c.Client(), action)

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
