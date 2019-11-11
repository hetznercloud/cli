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

	"github.com/cheggaaa/pb/v3"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var ErrConfigPathUnknown = errors.New("config file path unknown")

const (
	progressCircleTpl = `{{ cycle . "↖" "↗" "↘" "↙" }}`
	progressBarTpl    = `{{ etime . }} {{ bar . "" "=" }} {{ percent . }}`
)

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
		if c.Debug {
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
		progress := pb.New(100)
		progress.SetWidth(50) // width of progress bar is too large by default
		progress.SetTemplateString(progressBarTpl)
		progress.Start()
		defer progress.Finish()

		for {
			select {
			case err := <-errCh:
				if err == nil {
					progress.SetCurrent(100)
				}
				return err
			case p := <-progressCh:
				progress.SetCurrent(int64(p))
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
	const (
		done   = "done"
		failed = "failed"
	)

	for _, action := range actions {
		resources := make(map[string]int)
		for _, resource := range action.Resources {
			resources[string(resource.Type)] = resource.ID
		}

		var waitingFor string
		switch action.Command {
		default:
			waitingFor = fmt.Sprintf("Waiting for action %s to have finished... ", action.Command)
		case "start_server":
			waitingFor = fmt.Sprintf("Waiting for server %d to have started... ", resources["server"])
		case "attach_volume":
			waitingFor = fmt.Sprintf("Waiting for volume %d to have been attached to server %d... ", resources["volume"], resources["server"])
		}

		if c.Terminal() {
			progress := pb.New(1) // total progress of 1 will do since we use a circle here
			progress.SetTemplateString(waitingFor + progressCircleTpl)
			progress.Start()
			defer progress.Finish()

			_, errCh := c.Client().Action.WatchProgress(ctx, action)
			if err := <-errCh; err != nil {
				progress.SetTemplateString(waitingFor + failed)
				return err
			}
			progress.SetTemplateString(waitingFor + done)
		} else {
			fmt.Print(waitingFor)

			_, errCh := c.Client().Action.WatchProgress(ctx, action)
			if err := <-errCh; err != nil {
				fmt.Println(failed)
				return err
			}
			fmt.Println(done)
		}
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
