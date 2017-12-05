package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

type CLI struct {
	Token    string
	Endpoint string
	JSON     bool

	Config *Config

	RootCommand *cobra.Command

	client *hcloud.Client
}

func NewCLI() *CLI {
	cli := &CLI{}
	cli.RootCommand = NewRootCommand(cli)
	return cli
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
		c.client = hcloud.NewClient(
			hcloud.WithToken(c.Token),
			hcloud.WithEndpoint(c.Endpoint),
		)
	}
	return c.client
}

// Terminal returns whether the CLI is run in a terminal.
func (c *CLI) Terminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

const ESC = 27

// ClearLine clears the previous line.
func (c *CLI) ClearLine() {
	fmt.Printf("%c[%dA", ESC, 1) // move the cursor up
	fmt.Printf("%c[2K", ESC)     // clear the line
}
