package state

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/ui"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func (c *state) WaitForActions(cmd *cobra.Command, ctx context.Context, actions ...*hcloud.Action) error {
	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		return c.Client().Action().WaitFor(ctx, actions...)
	}

	return waitForActions(c.Client().Action(), ctx, actions...)
}

func waitForActions(client hcapi2.ActionClient, ctx context.Context, actions ...*hcloud.Action) (err error) {
	progressGroup := ui.NewProgressGroup(os.Stderr)
	progressByAction := make(map[int64]ui.Progress, len(actions))
	for _, action := range actions {
		progress := progressGroup.Add(waitForActionMessage(action))
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

func waitForActionMessage(action *hcloud.Action) string {
	resources := make([]string, 0, len(action.Resources))
	for _, resource := range action.Resources {
		resources = append(resources, fmt.Sprintf("%s: %d", resource.Type, resource.ID))
	}
	return fmt.Sprintf("Waiting for action %s to complete (%v)", action.Command, strings.Join(resources, ", "))
}
