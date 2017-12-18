package cli

import (
	"context"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

func yesno(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func na(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}

func waitAction(ctx context.Context, client *hcloud.Client, action *hcloud.Action) <-chan error {
	ch := make(chan error, 1)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				ch <- ctx.Err()
				return
			case <-ticker.C:
				break
			}

			action, _, err := client.Action.Get(ctx, action.ID)
			if err != nil {
				ch <- ctx.Err()
				return
			}

			switch action.Status {
			case hcloud.ActionStatusRunning:
				break
			case hcloud.ActionStatusSuccess:
				ch <- nil
				return
			case hcloud.ActionStatusError:
				ch <- action.Error()
				return
			}
		}
	}()

	return ch
}
