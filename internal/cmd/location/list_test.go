package location_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := location.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.Location.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.LocationListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Location{
			{
				ID:          1,
				Name:        "fsn1",
				NetworkZone: hcloud.NetworkZoneEUCentral,
				Country:     "DE",
				City:        "Falkenstein",
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID   NAME   DESCRIPTION   NETWORK ZONE   COUNTRY   CITY
1    fsn1   -             eu-central     DE        Falkenstein
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
