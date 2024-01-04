package firewall_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Firewall{ID: 123}, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		Update(gomock.Any(), &hcloud.Firewall{ID: 123}, hcloud.FirewallUpdateOpts{
			Name: "new-name",
		})

	out, _, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Firewall 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
