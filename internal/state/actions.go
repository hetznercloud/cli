package state

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/ui"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func (c *state) WaitForActions(ctx context.Context, _ *cobra.Command, actions ...*hcloud.Action) error {
	quiet, err := config.OptionQuiet.Get(c.Config())
	if err != nil {
		return err
	}
	if quiet {
		return c.Client().Action().WaitFor(ctx, actions...)
	}

	return waitForActions(ctx, c.Client().Action(), actions...)
}

func waitForActions(ctx context.Context, client hcapi2.ActionClient, actions ...*hcloud.Action) (err error) {
	progressGroup := ui.NewProgressGroup(os.Stderr)
	progressByAction := make(map[int64]ui.Progress, len(actions))
	for _, action := range actions {
		progress := progressGroup.Add(
			ui.ActionMessage(action),
			ui.ActionResourcesMessage(action.Resources...),
		)
		progressByAction[action.ID] = progress
	}

	if err = progressGroup.Start(); err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, progressGroup.Stop())
	}()

	return client.WaitForFunc(ctx, func(update *hcloud.Action) error {
		switch update.Status {
		case hcloud.ActionStatusRunning:
			progressByAction[update.ID].SetCurrent(update.Progress)
		case hcloud.ActionStatusSuccess:
			progressByAction[update.ID].SetSuccess()
		case hcloud.ActionStatusError:
			progressByAction[update.ID].SetError()
			return update.Error()
		}

		return nil
	}, actions...)
}
