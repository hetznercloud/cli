package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := updateCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer)
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
		Update(
			gomock.Any(),
			primaryip,
			hcloud.PrimaryIPUpdateOpts{
				Name: "foobar",
			},
		).
		Return(
			&hcloud.PrimaryIP{ID: 13, Name: "foobar"},
			&hcloud.Response{},
			nil,
		)

	out, err := fx.Run(cmd, []string{"13", "--name=foobar"})

	expOut := "Primary IP 13 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
