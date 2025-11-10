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

func TestAddRecords(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := rrset.AddRecordsCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
	}

	rrSet := &hcloud.ZoneRRSet{
		Zone: z,
		Name: "www",
		Type: hcloud.ZoneRRSetTypeA,
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		AddRRSetRecords(gomock.Any(), rrSet, hcloud.ZoneRRSetAddRecordsOpts{
			Records: []hcloud.ZoneRRSetRecord{{Value: "198.51.100.1"}},
			TTL:     nil,
		}).
		Return(&hcloud.Action{ID: 321}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

	out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "--record", "198.51.100.1"})

	expOut := "Added records on Zone RRSet www A\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
