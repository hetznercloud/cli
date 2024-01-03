package state

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const (
	progressCircleTpl = `{{ cycle . " .  " "  . " "   ." "  . " }}`
	progressBarTpl    = `{{ etime . }} {{ bar . "" "=" }} {{ percent . }}`
)

func Wrap(s State, f func(State, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(s, cmd, args)
	}
}

// StdoutIsTerminal returns whether the CLI is run in a terminal.
func StdoutIsTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func (c *state) ActionProgress(cmd *cobra.Command, ctx context.Context, action *hcloud.Action) error {
	return c.ActionsProgresses(cmd, ctx, []*hcloud.Action{action})
}

func (c *state) ActionsProgresses(cmd *cobra.Command, ctx context.Context, actions []*hcloud.Action) error {
	progressCh, errCh := c.hcloudClient.Action.WatchOverallProgress(ctx, actions)

	if StdoutIsTerminal() {
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

func (c *state) EnsureToken(_ *cobra.Command, _ []string) error {
	if c.Token == "" {
		return errors.New("no active context or token (see `hcloud context --help`)")
	}
	return nil
}

func (c *state) WaitForActions(cmd *cobra.Command, ctx context.Context, actions []*hcloud.Action) error {
	for _, action := range actions {
		resources := make(map[string]int64)
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

		_, errCh := c.hcloudClient.Action.WatchProgress(ctx, action)

		err := DisplayProgressCircle(cmd, errCh, waitingFor)
		if err != nil {
			return err
		}
	}

	return nil
}

func DisplayProgressCircle(cmd *cobra.Command, errCh <-chan error, waitingFor string) error {
	const (
		done     = "done"
		failed   = "failed"
		ellipsis = " ... "
	)

	if StdoutIsTerminal() {
		_, _ = fmt.Fprintln(os.Stderr, waitingFor)

		progress := pb.New(1) // total progress of 1 will do since we use a circle here
		progress.SetTemplateString(progressCircleTpl)
		progress.Start()
		defer progress.Finish()

		if err := <-errCh; err != nil {
			progress.SetTemplateString(ellipsis + failed)
			return err
		}
		progress.SetTemplateString(ellipsis + done)
	} else {
		_, _ = fmt.Fprint(os.Stderr, waitingFor+ellipsis)

		if err := <-errCh; err != nil {
			_, _ = fmt.Fprintln(os.Stderr, failed)
			return err
		}
		_, _ = fmt.Fprintln(os.Stderr, done)
	}
	return nil
}
