package screenshot

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestTakeScreenshot(t *testing.T) {
	t.Skip()

	ctx := t.Context()

	// TODO: Load token
	client := hcloud.NewClient(hcloud.WithToken("TOKEN"))

	// TODO: Create or use existing server
	result, _, err := client.Server.RequestConsole(ctx, &hcloud.Server{ID: 0})
	require.NoError(t, err)

	err = client.Action.WaitFor(ctx, result.Action)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	err = TakeScreenshot(ctx, result.WSSURL, "screenshot.png")
	require.NoError(t, err)
}
