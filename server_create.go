package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "create",
		Short:            "Create server",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE:             cli.wrap(runServerCreate),
	}
	cmd.Flags().String("name", "", "Server name")
	cmd.Flags().String("type", "", "Server type (id or name)")
	cmd.Flags().String("image", "", "Image (id or name)")
	return cmd
}

func runServerCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	name, _ := cmd.Flags().GetString("name")
	serverType, _ := cmd.Flags().GetString("type")
	image, _ := cmd.Flags().GetString("image")

	opts := hcloud.ServerCreateOpts{
		Name: name,
		ServerType: hcloud.ServerType{
			Name: serverType,
		},
		Image: hcloud.Image{
			Name: image,
		},
	}

	result, _, err := cli.Client().Server.Create(ctx, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Server %d created\n", result.Server.ID)
	initStarted := time.Now()
	terminal := cli.Terminal()

	if terminal {
		fmt.Println("Initializing server ...")
	}

progress:
	for {
		action, _, err := cli.Client().Action.Get(ctx, result.Action.ID)
		if err != nil {
			return err
		}
		switch action.Status {
		case "running":
			if terminal {
				cli.ClearLine()
				fmt.Printf("Initializing server ... %d%%\n", action.Progress)
			}
			time.Sleep(100 * time.Millisecond)
		case "success":
			duration := time.Now().Sub(initStarted).Round(time.Millisecond)
			if terminal {
				cli.ClearLine()
			}
			fmt.Printf("Server initialized in %s\n", duration)
			break progress
		case "error":
			fmt.Println(action.ErrorMessage)
			break progress
		}
	}

	return nil
}
