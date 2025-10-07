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

func TestSetRecords(t *testing.T) {

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

	records := []hcloud.ZoneRRSetRecord{{Value: "198.51.100.1"}, {Value: "198.51.100.2"}}

	t.Run("default", func(t *testing.T) {
		fx := testutil.NewFixture(t)
		defer fx.Finish()

		cmd := rrset.SetRecordsCmd.CobraCommand(fx.State())
		fx.ExpectEnsureToken()

		fx.Client.ZoneClient.EXPECT().
			Get(gomock.Any(), "example.com").
			Return(z, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
			Return(rrSet, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			SetRRSetRecords(gomock.Any(), rrSet, hcloud.ZoneRRSetSetRecordsOpts{
				Records: records,
			}).
			Return(&hcloud.Action{ID: 321}, nil, nil)
		fx.ActionWaiter.EXPECT().
			WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

		out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "--record", "198.51.100.1", "--record", "198.51.100.2"})

		expOut := "Set records on Zone RRSet www A\n"

		require.NoError(t, err)
		assert.Equal(t, ExperimentalWarning, errOut)
		assert.Equal(t, expOut, out)
	})

	t.Run("non-existing", func(t *testing.T) {
		fx := testutil.NewFixture(t)
		defer fx.Finish()

		cmd := rrset.SetRecordsCmd.CobraCommand(fx.State())
		fx.ExpectEnsureToken()

		fx.Client.ZoneClient.EXPECT().
			Get(gomock.Any(), "example.com").
			Return(z, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
			Return(nil, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			CreateRRSet(gomock.Any(), z, hcloud.ZoneRRSetCreateOpts{
				Name:    rrSet.Name,
				Type:    rrSet.Type,
				Records: records,
			}).
			Return(hcloud.ZoneRRSetCreateResult{
				RRSet:  rrSet,
				Action: &hcloud.Action{ID: 321},
			}, nil, nil)
		fx.ActionWaiter.EXPECT().
			WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

		out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "--record", "198.51.100.1", "--record", "198.51.100.2"})

		expOut := "Created and set records on Zone RRSet www A\n"

		require.NoError(t, err)
		assert.Equal(t, ExperimentalWarning, errOut)
		assert.Equal(t, expOut, out)
	})

	t.Run("empty-records", func(t *testing.T) {
		fx := testutil.NewFixture(t)
		defer fx.Finish()

		cmd := rrset.SetRecordsCmd.CobraCommand(fx.State())
		fx.ExpectEnsureToken()

		fx.Client.ZoneClient.EXPECT().
			Get(gomock.Any(), "example.com").
			Return(z, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
			Return(rrSet, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			DeleteRRSet(gomock.Any(), rrSet).
			Return(hcloud.ZoneRRSetDeleteResult{Action: &hcloud.Action{ID: 321}}, nil, nil)
		fx.ActionWaiter.EXPECT().
			WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})

		out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "--record="})

		expOut := "Zone RRSet www A deleted\n"

		require.NoError(t, err)
		assert.Equal(t, ExperimentalWarning, errOut)
		assert.Equal(t, expOut, out)
	})

	t.Run("empty-records-non-existing", func(t *testing.T) {
		fx := testutil.NewFixture(t)
		defer fx.Finish()

		cmd := rrset.SetRecordsCmd.CobraCommand(fx.State())
		fx.ExpectEnsureToken()

		fx.Client.ZoneClient.EXPECT().
			Get(gomock.Any(), "example.com").
			Return(z, nil, nil)
		fx.Client.ZoneClient.EXPECT().
			GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
			Return(nil, nil, nil)

		out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "--record="})

		expOut := "Zone RRSet www A doesn't exist. No action necessary.\n"

		require.NoError(t, err)
		assert.Equal(t, ExperimentalWarning, errOut)
		assert.Equal(t, expOut, out)
	})
}
