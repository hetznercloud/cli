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
