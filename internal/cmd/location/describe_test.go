package location_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := location.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LocationClient.EXPECT().
		Get(gomock.Any(), "hel1").
		Return(&hcloud.Location{
			ID:          3,
			Name:        "hel1",
			Description: "Helsinki DC Park 1",
			NetworkZone: "eu-central",
			Country:     "FI",
			City:        "Helsinki",
			Latitude:    60.169855,
			Longitude:   24.938379,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"hel1"})

	expOut := `ID:            3
Name:          hel1
Description:   Helsinki DC Park 1
Network Zone:  eu-central
Country:       FI
City:          Helsinki
Latitude:      60.169855
Longitude:     24.938379
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
