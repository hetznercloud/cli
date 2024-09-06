package datacenter_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/datacenter"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := datacenter.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.DatacenterClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.DatacenterListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Datacenter{
			{
				ID:          4,
				Name:        "fsn1-dc14",
				Location:    &hcloud.Location{Name: "fsn1"},
				Description: "Falkenstein 1 virtual DC 14",
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID   NAME        DESCRIPTION                   LOCATION
4    fsn1-dc14   Falkenstein 1 virtual DC 14   fsn1
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestListJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := datacenter.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.DatacenterClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.DatacenterListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Datacenter{
			{
				ID:          4,
				Name:        "fsn1-dc14",
				Location:    &hcloud.Location{Name: "fsn1"},
				Description: "Falkenstein 1 virtual DC 14",
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"-o=json"})

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.JSONEq(t, `
[
  {
    "id": 4,
    "name": "fsn1-dc14",
    "description": "Falkenstein 1 virtual DC 14",
    "location": {
      "id": 0,
      "name": "fsn1",
      "description": "",
      "country": "",
      "city": "",
      "latitude": 0,
      "longitude": 0,
      "network_zone": ""
    },
    "server_types": {
      "supported": null,
      "available_for_migration": null,
      "available": null
    }
  }
]`, out)
}
