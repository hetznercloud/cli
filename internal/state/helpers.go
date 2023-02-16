package state

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/hetznercloud/cli/internal/hcapi"
	"github.com/hetznercloud/cli/internal/version"
)

const (
	progressCircleTpl = `{{ cycle . " .  " "  . " "   ." "  . " }}`
	progressBarTpl    = `{{ etime . }} {{ bar . "" "=" }} {{ percent . }}`
)

func (c *State) Wrap(f func(*State, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(c, cmd, args)
	}
}

func (c *State) Client() *hcloud.Client {
	if c.client == nil {
		opts := []hcloud.ClientOption{
			hcloud.WithToken(c.Token),
			hcloud.WithApplication("hcloud-cli", version.Version),
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
		// TODO Somehow pass here
		// pollInterval, _ := c.RootCommand.PersistentFlags().GetDuration("poll-interval")
		pollInterval := 500 * time.Millisecond
		if pollInterval > 0 {
			opts = append(opts, hcloud.WithPollInterval(pollInterval))
		}
		c.client = hcloud.NewClient(opts...)
	}
	return c.client
}

func (c *State) FirewallNames() []string {
	if c.firewallClient == nil {
		client := c.Client()
		c.firewallClient = &hcapi.FirewallClient{FirewallClient: &client.Firewall}
	}
	return c.firewallClient.FirewallNames()
}

func (c *State) FirewallLabelKeys(idOrName string) []string {
	if c.firewallClient == nil {
		client := c.Client()
		c.firewallClient = &hcapi.FirewallClient{FirewallClient: &client.Firewall}
	}
	return c.firewallClient.FirewallLabelKeys(idOrName)
}

func (c *State) FloatingIPNames() []string {
	if c.floatingIPClient == nil {
		client := c.Client()
		c.floatingIPClient = &hcapi.FloatingIPClient{FloatingIPClient: &client.FloatingIP}
	}
	return c.floatingIPClient.FloatingIPNames()
}

func (c *State) FloatingIPLabelKeys(idOrName string) []string {
	if c.floatingIPClient == nil {
		client := c.Client()
		c.floatingIPClient = &hcapi.FloatingIPClient{FloatingIPClient: &client.FloatingIP}
	}
	return c.floatingIPClient.FloatingIPLabelKeys(idOrName)
}

func (c *State) ISONames() []string {
	if c.isoClient == nil {
		client := c.Client()
		c.isoClient = &hcapi.ISOClient{ISOClient: &client.ISO}
	}
	return c.isoClient.ISONames()
}

func (c *State) ImageNames() []string {
	if c.imageClient == nil {
		client := c.Client()
		c.imageClient = &hcapi.ImageClient{ImageClient: &client.Image}
	}
	return c.imageClient.ImageNames()
}

func (c *State) ImageLabelKeys(idOrName string) []string {
	if c.imageClient == nil {
		client := c.Client()
		c.imageClient = &hcapi.ImageClient{ImageClient: &client.Image}
	}
	return c.imageClient.ImageLabelKeys(idOrName)
}

func (c *State) LocationNames() []string {
	if c.locationClient == nil {
		client := c.Client()
		c.locationClient = &hcapi.LocationClient{LocationClient: &client.Location}
	}
	return c.locationClient.LocationNames()
}

func (c *State) NetworkZoneNames() []string {
	if c.locationClient == nil {
		client := c.Client()
		c.locationClient = &hcapi.LocationClient{LocationClient: &client.Location}
	}
	return c.locationClient.NetworkZoneNames()
}

func (c *State) DataCenterNames() []string {
	if c.dataCenterClient == nil {
		client := c.Client()
		c.dataCenterClient = &hcapi.DataCenterClient{DatacenterClient: &client.Datacenter}
	}
	return c.dataCenterClient.DataCenterNames()
}

func (c *State) SSHKeyNames() []string {
	if c.sshKeyClient == nil {
		client := c.Client()
		c.sshKeyClient = &hcapi.SSHKeyClient{SSHKeyClient: &client.SSHKey}
	}
	return c.sshKeyClient.SSHKeyNames()
}

func (c *State) SSHKeyLabelKeys(idOrName string) []string {
	if c.sshKeyClient == nil {
		client := c.Client()
		c.sshKeyClient = &hcapi.SSHKeyClient{SSHKeyClient: &client.SSHKey}
	}
	return c.sshKeyClient.SSHKeyLabelKeys(idOrName)
}

func (c *State) VolumeNames() []string {
	if c.volumeClient == nil {
		client := c.Client()
		c.volumeClient = &hcapi.VolumeClient{VolumeClient: &client.Volume}
	}
	return c.volumeClient.VolumeNames()
}

func (c *State) VolumeLabelKeys(idOrName string) []string {
	if c.volumeClient == nil {
		client := c.Client()
		c.volumeClient = &hcapi.VolumeClient{VolumeClient: &client.Volume}
	}
	return c.volumeClient.VolumeLabelKeys(idOrName)
}

// Terminal returns whether the CLI is run in a terminal.
func (c *State) Terminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func (c *State) ActionProgress(ctx context.Context, action *hcloud.Action) error {
	return c.ActionsProgresses(ctx, []*hcloud.Action{action})
}

func (c *State) ActionsProgresses(ctx context.Context, actions []*hcloud.Action) error {
	progressCh, errCh := c.Client().Action.WatchOverallProgress(ctx, actions)

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

func (c *State) EnsureToken(cmd *cobra.Command, args []string) error {
	if c.Token == "" {
		return errors.New("no active context or token (see `hcloud context --help`)")
	}
	return nil
}

func (c *State) WaitForActions(ctx context.Context, actions []*hcloud.Action) error {
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

func (c *State) ServerTypeNames() []string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerTypeNames()
}

func (c *State) ServerNames() []string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerNames()
}

func (c *State) ServerLabelKeys(idOrName string) []string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerLabelKeys(idOrName)
}

func (c *State) ServerName(id int) string {
	if c.serverClient == nil {
		client := c.Client()
		c.serverClient = &hcapi.ServerClient{
			ServerClient: &client.Server,
			ServerTypes:  &client.ServerType,
		}
	}
	return c.serverClient.ServerName(id)
}

func (c *State) NetworkNames() []string {
	if c.networkClient == nil {
		client := c.Client()
		c.networkClient = &hcapi.NetworkClient{NetworkClient: &client.Network}
	}
	return c.networkClient.NetworkNames()
}

func (c *State) NetworkName(id int) string {
	if c.networkClient == nil {
		client := c.Client()
		c.networkClient = &hcapi.NetworkClient{NetworkClient: &client.Network}
	}
	return c.networkClient.NetworkName(id)
}

func (c *State) NetworkLabelKeys(idOrName string) []string {
	if c.networkClient == nil {
		client := c.Client()
		c.networkClient = &hcapi.NetworkClient{NetworkClient: &client.Network}
	}
	return c.networkClient.NetworkLabelKeys(idOrName)
}

func (c *State) LoadBalancerNames() []string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerNames()
}

func (c *State) LoadBalancerName(id int) string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerName(id)
}

func (c *State) LoadBalancerLabelKeys(idOrName string) []string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerLabelKeys(idOrName)
}

func (c *State) LoadBalancerTypeNames() []string {
	if c.loadBalancerClient == nil {
		client := c.Client()
		c.loadBalancerClient = &hcapi.LoadBalancerClient{
			LoadBalancerClient: &client.LoadBalancer,
			TypeClient:         &client.LoadBalancerType,
		}
	}
	return c.loadBalancerClient.LoadBalancerTypeNames()
}

func (c *State) PlacementGroupNames() []string {
	if c.placementGroupClient == nil {
		client := c.Client()
		c.placementGroupClient = &hcapi.PlacementGroupClient{PlacementGroupClient: &client.PlacementGroup}
	}
	return c.placementGroupClient.PlacementGroupNames()
}
