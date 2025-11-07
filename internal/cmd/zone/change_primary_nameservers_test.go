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

func TestChangePrimaryNameservers(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.ChangePrimaryNameserversCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{ID: 123, Name: "example.com"}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		ChangePrimaryNameservers(gomock.Any(), z, hcloud.ZoneChangePrimaryNameserversOpts{
			PrimaryNameservers: []hcloud.ZoneChangePrimaryNameserversOptsPrimaryNameserver{
				{Address: "198.51.100.1", Port: 53},
				{Address: "203.0.113.1"},
				{Address: "203.0.113.1", Port: 53, TSIGKey: "example-key", TSIGAlgorithm: hcloud.ZoneTSIGAlgorithmHMACSHA256},
			},
		}).
		Return(&hcloud.Action{ID: 321}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, errOut, err := fx.Run(cmd, []string{"example.com", "--primary-nameservers-file", "testdata/primary_nameservers.json"})

	expOut := "Primary nameservers for Zone example.com updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
