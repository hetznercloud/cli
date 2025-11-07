package zone_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeTTL(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.ChangeTTLCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{ID: 123, Name: "example.com"}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		ChangeTTL(gomock.Any(), z, hcloud.ZoneChangeTTLOpts{TTL: 1337}).
		Return(&hcloud.Action{ID: 321}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, errOut, err := fx.Run(cmd, []string{"--ttl", "1337", "123"})

	expOut := "Changed default TTL on Zone example.com\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
