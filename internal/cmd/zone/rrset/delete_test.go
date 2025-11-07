package rrset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone/rrset"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := rrset.DeleteCmd.CobraCommand(fx.State())
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
		DeleteRRSet(gomock.Any(), rrSet).
		Return(hcloud.ZoneRRSetDeleteResult{Action: &hcloud.Action{ID: 42}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 42}}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A"})

	expOut := "Zone RRSet www A deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
