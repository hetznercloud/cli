package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer, fx.ActionWaiter)
	primaryip := &hcloud.PrimaryIP{ID: 13}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			primaryip,
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		Delete(
			gomock.Any(),
			primaryip,
		).
		Return(
			&hcloud.Response{},
			nil,
		)

	out, err := fx.Run(cmd, []string{"13"})

	expOut := "Primary IP 13 deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
