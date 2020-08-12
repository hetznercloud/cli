package cli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/hetznercloud/cli/internal/hcapi"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var ErrConfigPathUnknown = errors.New("config file path unknown")

const (
	progressCircleTpl = `{{ cycle . " .  " "  . " "   ." "  . " }}`
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

	client             *hcloud.Client
	isoClient          *hcapi.ISOClient
	imageClient        *hcapi.ImageClient
	locationClient     *hcapi.LocationClient
	dataCenterClient   *hcapi.DataCenterClient
	sshKeyClient       *hcapi.SSHKeyClient
	volumeClient       *hcapi.VolumeClient
	certificateClient  *hcapi.CertificateClient
	floatingIPClient   *hcapi.FloatingIPClient
	networkClient      *hcapi.NetworkClient
	loadBalancerClient *hcapi.LoadBalancerClient
	serverClient       *hcapi.ServerClient
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

func (c *CLI) CertificateNames() []string {
	if c.certificateClient == nil {
		client := c.Client()
		c.certificateClient = &hcapi.CertificateClient{CertificateClient: &client.Certificate}
	}
	return c.certificateClient.CertificateNames()
}

func (c *CLI) CertificateLabelKeys(idOrName string) []string {
	if c.certificateClient == nil {
		client := c.Client()
		c.certificateClient = &hcapi.CertificateClient{CertificateClient: &client.Certificate}
	}
	return c.certificateClient.CertificateLabelKeys(idOrName)
}

func (c *CLI) FloatingIPNames() []string {
	if c.floatingIPClient == nil {
		client := c.Client()
		c.floatingIPClient = &hcapi.FloatingIPClient{FloatingIPClient: &client.FloatingIP}
	}
	return c.floatingIPClient.FloatingIPNames()
}

func (c *CLI) FloatingIPLabelKeys(idOrName string) []string {
	if c.floatingIPClient == nil {
		client := c.Client()
		c.floatingIPClient = &hcapi.FloatingIPClient{FloatingIPClient: &client.FloatingIP}
	}
	return c.floatingIPClient.FloatingIPLabelKeys(idOrName)
}

func (c *CLI) ISONames() []string {
	if c.isoClient == nil {
		client := c.Client()
		c.isoClient = &hcapi.ISOClient{ISOClient: &client.ISO}
	}
	return c.isoClient.ISONames()
}

func (c *CLI) ImageNames() []string {
	if c.imageClient == nil {
		client := c.Client()
		c.imageClient = &hcapi.ImageClient{ImageClient: &client.Image}
	}
	return c.isoClient.ISONames()
}

func (c *CLI) ImageLabelKeys(idOrName string) []string {
	if c.imageClient == nil {
		client := c.Client()
		c.imageClient = &hcapi.ImageClient{ImageClient: &client.Image}
	}
	return c.imageClient.ImageLabelKeys(idOrName)
}

func (c *CLI) LocationNames() []string {
	if c.locationClient == nil {
		client := c.Client()
		c.locationClient = &hcapi.LocationClient{LocationClient: &client.Location}
	}
	return c.locationClient.LocationNames()
}

func (c *CLI) DataCenterNames() []string {
	if c.dataCenterClient == nil {
		client := c.Client()
		c.dataCenterClient = &hcapi.DataCenterClient{DatacenterClient: &client.Datacenter}
	}
	return c.dataCenterClient.DataCenterNames()
}

func (c *CLI) SSHKeyNames() []string {
	if c.sshKeyClient == nil {
		client := c.Client()
		c.sshKeyClient = &hcapi.SSHKeyClient{SSHKeyClient: &client.SSHKey}
	}
	return c.sshKeyClient.SSHKeyNames()
}

func (c *CLI) SSHKeyLabelKeys(idOrName string) []string {
	if c.sshKeyClient == nil {
		client := c.Client()
		c.sshKeyClient = &hcapi.SSHKeyClient{SSHKeyClient: &client.SSHKey}
	}
	return c.sshKeyClient.SSHKeyLabelKeys(idOrName)
}

func (c *CLI) VolumeNames() []string {
	if c.volumeClient == nil {
		client := c.Client()
		c.volumeClient = &hcapi.VolumeClient{VolumeClient: &client.Volume}
	}
	return c.volumeClient.VolumeNames()
}

func (c *CLI) VolumeLabelKeys(idOrName string) []string {
	if c.volumeClient == nil {
		client := c.Client()
		c.volumeClient = &hcapi.VolumeClient{VolumeClient: &client.Volume}
	}
	return c.volumeClient.VolumeLabelKeys(idOrName)
}

// Terminal returns whether the CLI is run in a terminal.
func (c *CLI) Terminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func (c *CLI) ActionProgress(ctx context.Context, action *hcloud.Action) error {
	progressCh, errCh := c.Client().Action.WatchProgress(ctx, action)

	if c.Terminal() {
		progress := pb.New(100)
		progress.SetMaxWidth(50) // width of progress bar is too large by default
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
		done     = "done"
		failed   = "failed"
		ellipsis = " ... "
	)

	for _, action := range actions {
		resources := make(map[string]int)
		for _, resource := range action.Resources {
			resources[string(resource.Type)] = resource.ID
		}

		var waitingFor string
		switch action.Command {
		default:
			waitingFor = fmt.Sprintf("Waiting for action %s to have finished", action.Command)
		case "start_server":
			waitingFor = fmt.Sprintf("Waiting for server %d to have started", resources["server"])
		case "attach_volume":
			waitingFor = fmt.Sprintf("Waiting for volume %d to have been attached to server %d", resources["volume"], resources["server"])
		}

		if c.Terminal() {
			fmt.Println(waitingFor)
			progress := pb.New(1) // total progress of 1 will do since we use a circle here
			progress.SetTemplateString(progressCircleTpl)
			progress.Start()
			defer progress.Finish()

			_, errCh := c.Client().Action.WatchProgress(ctx, action)
			if err := <-errCh; err != nil {
				progress.SetTemplateString(ellipsis + failed)
				return err
			}
			progress.SetTemplateString(ellipsis + done)
		} else {
			fmt.Print(waitingFor + ellipsis)

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

func (c *CLI) ServerTypeNames() []string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerTypeNames()
}

func (c *CLI) ServerNames() []string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerNames()
}

func (c *CLI) ServerLabelKeys(idOrName string) []string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerLabelKeys(idOrName)
}

func (c *CLI) ServerName(id int) string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerName(id)
}

func (c *CLI) NetworkNames() []string {
	if c.networkClient == nil {
		client := c.Client()
		c.networkClient = &hcapi.NetworkClient{NetworkClient: &client.Network}
	}
	return c.networkClient.NetworkNames()
}

func (c *CLI) NetworkName(id int) string {
	if c.networkClient == nil {
		client := c.Client()
		c.networkClient = &hcapi.NetworkClient{NetworkClient: &client.Network}
	}
	return c.networkClient.NetworkName(id)
}

func (c *CLI) NetworkLabelKeys(idOrName string) []string {
	if c.networkClient == nil {
		client := c.Client()
		c.networkClient = &hcapi.NetworkClient{NetworkClient: &client.Network}
	}
	return c.networkClient.NetworkLabelKeys(idOrName)
}

func (c *CLI) LoadBalancerNames() []string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerNames()
}

func (c *CLI) LoadBalancerLabelKeys(idOrName string) []string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerLabelKeys(idOrName)
}

func (c *CLI) LoadBalancerTypeNames() []string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerTypeNames()
}
