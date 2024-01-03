package firewall

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	firewall := &hcloud.Firewall{
		ID:   123,
		Name: "test",
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(firewall, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		Delete(gomock.Any(), firewall).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "firewall test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
