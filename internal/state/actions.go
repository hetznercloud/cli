package state

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/ui"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func (c *state) WaitForActions(ctx context.Context, cmd *cobra.Command, actions ...*hcloud.Action) error {
	return c.waitForActions(ctx, cmd, actions...)
}

func (c *state) waitForActions(ctx context.Context, cmd *cobra.Command, actions ...*hcloud.Action) error {
	quiet, err := config.OptionQuiet.Get(c.Config())
	if err != nil {
		return err
	}
	if quiet {
		return waitForActionsQuiet(ctx, c.Client().Action(), actions...)
	}

	out := c.stderr
	if cmd != nil {
		out = cmd.ErrOrStderr()
	}
	return waitForActions(ctx, c.Client().Action(), out, actions...)
}

type ActionFailure struct {
	ActionID int64
	Err      error
}

func (e ActionFailure) Error() string {
	return fmt.Sprintf("action %d failed: %v", e.ActionID, e.Err)
}

func (e ActionFailure) Unwrap() error {
	return e.Err
}

type ActionWaitError struct {
	Failures []ActionFailure
	Cause    error
}

func (e *ActionWaitError) Error() string {
	errs := make([]error, 0, len(e.Failures)+1)
	if e.Cause != nil {
		errs = append(errs, e.Cause)
	}
	for _, failure := range e.Failures {
		errs = append(errs, failure)
	}
	return errors.Join(errs...).Error()
}

func (e *ActionWaitError) Unwrap() []error {
	errs := make([]error, 0, len(e.Failures)+1)
	if e.Cause != nil {
		errs = append(errs, e.Cause)
	}
	for _, failure := range e.Failures {
		errs = append(errs, failure)
	}
	return errs
}

func waitForActionsQuiet(ctx context.Context, client hcapi2.ActionClient, actions ...*hcloud.Action) error {
	return collectActionFailures(ctx, client, nil, actions...)
}

func waitForActions(ctx context.Context, client hcapi2.ActionClient, out io.Writer, actions ...*hcloud.Action) (err error) {
	progressGroup := ui.NewProgressGroup(out)
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

	return collectActionFailures(ctx, client, func(update *hcloud.Action) {
		progress, ok := progressByAction[update.ID]
		if !ok {
			return
		}
		switch update.Status {
		case hcloud.ActionStatusRunning:
			progress.SetCurrent(update.Progress)
		case hcloud.ActionStatusSuccess:
			progress.SetSuccess()
		case hcloud.ActionStatusError:
			progress.SetError()
		}
	}, actions...)
}

func collectActionFailures(
	ctx context.Context,
	client hcapi2.ActionClient,
	onUpdate func(*hcloud.Action),
	actions ...*hcloud.Action,
) error {
	failures := make(map[int64]ActionFailure)
	err := client.WaitForFunc(ctx, func(update *hcloud.Action) error {
		if onUpdate != nil {
			onUpdate(update)
		}
		if update.Status == hcloud.ActionStatusError {
			failures[update.ID] = ActionFailure{ActionID: update.ID, Err: update.Error()}
		}
		return nil
	}, actions...)
	if err == nil && len(failures) == 0 {
		return nil
	}

	result := &ActionWaitError{Cause: err}
	for _, action := range actions {
		if failure, ok := failures[action.ID]; ok {
			result.Failures = append(result.Failures, failure)
			delete(failures, action.ID)
		}
	}
	for _, failure := range failures {
		result.Failures = append(result.Failures, failure)
	}
	return result
}
