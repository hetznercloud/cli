package floatingip

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

	floatingIP := &hcloud.FloatingIP{
		ID:   123,
		Name: "test",
	}

	fx.Client.FloatingIPClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(floatingIP, nil, nil)
	fx.Client.FloatingIPClient.EXPECT().
		Delete(gomock.Any(), floatingIP).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "Floating IP test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
