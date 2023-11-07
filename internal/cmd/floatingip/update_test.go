package floatingip

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

	cmd := UpdateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.FloatingIP{ID: 123}, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		Update(gomock.Any(), &hcloud.FloatingIP{ID: 123}, hcloud.FloatingIPUpdateOpts{
			Name: "new-name",
		})

	out, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Floating IP 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestUpdateDescription(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UpdateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.FloatingIP{ID: 123}, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		Update(gomock.Any(), &hcloud.FloatingIP{ID: 123}, hcloud.FloatingIPUpdateOpts{
			Description: "new-description",
		})

	out, err := fx.Run(cmd, []string{"123", "--description", "new-description"})

	expOut := "Floating IP 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
