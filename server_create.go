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
		Short:                 "Create server",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerCreate),
	}
	cmd.Flags().String("name", "", "Server name")

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
	cmd.Flags().String("datacenter", "", "Datacenter (ID or name)")
	cmd.Flags().IntSlice("ssh-key", nil, "ID of SSH key to inject (can be specified multiple times)")
	return cmd
}

func runServerCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	opts := optsFromFlags(cmd.Flags())
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

func optsFromFlags(flags *pflag.FlagSet) hcloud.ServerCreateOpts {
	name, _ := flags.GetString("name")
	serverType, _ := flags.GetString("type")
	image, _ := flags.GetString("image")
	location, _ := flags.GetString("location")
	datacenter, _ := flags.GetString("datacenter")
	sshKeys, _ := flags.GetIntSlice("ssh-key")

	opts := hcloud.ServerCreateOpts{
		Name: name,
		ServerType: hcloud.ServerType{
			Name: serverType,
		},
		Image: hcloud.Image{
			Name: image,
		},
	}
	for _, sshKey := range sshKeys {
		opts.SSHKeys = append(opts.SSHKeys, &hcloud.SSHKey{ID: sshKey})
	}
	if datacenter != "" {
		opts.Datacenter = &hcloud.Datacenter{Name: datacenter}
	}
	if location != "" {
		opts.Location = &hcloud.Location{Name: location}
	}

	return opts
}
