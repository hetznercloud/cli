package primaryip

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := updateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.PrimaryIPClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.PrimaryIP{ID: 123}, nil, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		Update(gomock.Any(), &hcloud.PrimaryIP{ID: 123}, hcloud.PrimaryIPUpdateOpts{
			Name:       "new-name",
			AutoDelete: nil,
		})

	out, _, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Primary IP 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestUpdateAutoDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := updateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.PrimaryIPClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.PrimaryIP{ID: 123}, nil, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		Update(gomock.Any(), &hcloud.PrimaryIP{ID: 123}, hcloud.PrimaryIPUpdateOpts{
			AutoDelete: hcloud.Ptr(false),
		})

	out, _, err := fx.Run(cmd, []string{"123", "--auto-delete=false"})

	expOut := "Primary IP 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
