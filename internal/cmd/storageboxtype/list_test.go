package storageboxtype_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storageboxtype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := storageboxtype.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.StorageBoxTypeClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.StorageBoxTypeListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
			},
		).
		Return([]*hcloud.StorageBoxType{
			{
				ID:                     42,
				Name:                   "bx11",
				Description:            "BX11",
				SnapshotLimit:          hcloud.Ptr(10),
				AutomaticSnapshotLimit: hcloud.Ptr(10),
				SubaccountsLimit:       200,
				Size:                   1073741824,
				Pricings:               []hcloud.StorageBoxTypeLocationPricing{},
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID   NAME   DESCRIPTION   SIZE      SNAPSHOT LIMIT   AUTOMATIC SNAPSHOT LIMIT   SUBACCOUNTS LIMIT
42   bx11   BX11          1.0 GiB   10               10                         200              
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestListColumnDeprecated(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := storageboxtype.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.StorageBoxTypeClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.StorageBoxTypeListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
			},
		).
		Return([]*hcloud.StorageBoxType{
			{
				ID:   123,
				Name: "deprecated",
				DeprecatableResource: hcloud.DeprecatableResource{
					Deprecation: &hcloud.DeprecationInfo{
						Announced:        time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
						UnavailableAfter: time.Date(2037, 8, 20, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			{
				ID:   124,
				Name: "current",
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"-o=columns=id,name,deprecated"})

	expOut := `ID    NAME         DEPRECATED                  
123   deprecated   Thu Aug 20 12:00:00 UTC 2037
124   current      -                           
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
