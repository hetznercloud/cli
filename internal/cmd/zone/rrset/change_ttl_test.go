package rrset_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone/rrset"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestChangeTTL(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := rrset.ChangeTTLCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
	}

	rrSet := &hcloud.ZoneRRSet{
		Zone: z,
		ID:   "www/A",
		Name: "www",
		Type: hcloud.ZoneRRSetTypeA,
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
		Return(rrSet, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		ChangeRRSetTTL(gomock.Any(), rrSet, hcloud.ZoneRRSetChangeTTLOpts{TTL: hcloud.Ptr(1337)}).
		Return(&hcloud.Action{ID: 321}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, errOut, err := fx.Run(cmd, []string{"--ttl", "1337", "example.com", "www", "A"})

	expOut := "Changed TTL on Zone RRSet www A\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
