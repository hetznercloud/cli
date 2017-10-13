package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const Endpoint = "https://api.hetzner.cloud/v1"

type CLI struct {
	Token    string
	Endpoint string
	JSON     bool

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

	if config.Token != "" {
		c.Token = config.Token
	}
	if config.Endpoint != "" {
		c.Endpoint = config.Endpoint
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
