package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newServerCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create FLAGS",
		Short:                 "Create a server",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerCreate),
	}
	cmd.Flags().String("name", "", "Server name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("type", "", "Server type (id or name)")
	cmd.Flag("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_servertype_names"},
	}
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("image", "", "Image (id or name)")
	cmd.Flag("image").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_image_names"},
	}
	cmd.MarkFlagRequired("image")

	cmd.Flags().String("location", "", "Location (ID or name)")
	cmd.Flag("location").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_location_names"},
	}

	cmd.Flags().String("datacenter", "", "Datacenter (ID or name)")
	cmd.Flag("datacenter").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_datacenter_names"},
	}

	cmd.Flags().StringSlice("ssh-key", nil, "ID or name of SSH key to inject (can be specified multiple times)")
	cmd.Flag("ssh-key").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_sshkey_names"},
	}
	return cmd
}

func runServerCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	opts, err := optsFromFlags(cli, cmd.Flags())
	if err != nil {
		return err
	}

	result, _, err := cli.Client().Server.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	if result.RootPassword != "" {
		fmt.Printf("Server %d created with root password: %s\n", result.Server.ID, result.RootPassword)
	} else {
		fmt.Printf("Server %d created\n", result.Server.ID)
	}

	return nil
}

func optsFromFlags(cli *CLI, flags *pflag.FlagSet) (opts hcloud.ServerCreateOpts, err error) {
	name, _ := flags.GetString("name")
	serverType, _ := flags.GetString("type")
	image, _ := flags.GetString("image")
	location, _ := flags.GetString("location")
	datacenter, _ := flags.GetString("datacenter")
	sshKeys, _ := flags.GetStringSlice("ssh-key")

	opts = hcloud.ServerCreateOpts{
		Name: name,
		ServerType: &hcloud.ServerType{
			Name: serverType,
		},
		Image: &hcloud.Image{
			Name: image,
		},
	}
	for _, sshKeyIDOrName := range sshKeys {
		var sshKey *hcloud.SSHKey
		sshKey, _, err = cli.Client().SSHKey.Get(cli.Context, sshKeyIDOrName)
		if err != nil {
			return
		}
		if sshKey == nil {
			err = fmt.Errorf("SSH key not found: %s", sshKeyIDOrName)
			return
		}
		opts.SSHKeys = append(opts.SSHKeys, sshKey)
	}
	if datacenter != "" {
		opts.Datacenter = &hcloud.Datacenter{Name: datacenter}
	}
	if location != "" {
		opts.Location = &hcloud.Location{Name: location}
	}

	return
}
