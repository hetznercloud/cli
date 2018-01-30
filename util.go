package cli

import (
	"context"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func yesno(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func na(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func datetime(t time.Time) string {
	return t.Local().Format(time.UnixDate)
}

func waitAction(ctx context.Context, client *hcloud.Client, action *hcloud.Action) (<-chan error, <-chan int) {
	errCh := make(chan error, 1)
	progressCh := make(chan int)

	go func() {
		defer close(errCh)
		defer close(progressCh)

		ticker := time.NewTicker(100 * time.Millisecond)

		sendProgress := func(p int) {
			select {
			case progressCh <- p:
				break
			default:
				break
			}
		}

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case <-ticker.C:
				break
			}

			action, _, err := client.Action.GetByID(ctx, action.ID)
			if err != nil {
				errCh <- ctx.Err()
				return
			}

			switch action.Status {
			case hcloud.ActionStatusRunning:
				sendProgress(action.Progress)
				break
			case hcloud.ActionStatusSuccess:
				sendProgress(100)
				errCh <- nil
				return
			case hcloud.ActionStatusError:
				errCh <- action.Error()
				return
			}
		}
	}()

	return errCh, progressCh
}

func chainRunE(fns ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range fns {
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
