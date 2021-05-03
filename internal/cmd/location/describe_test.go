package location_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := location.DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
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

	out, err := fx.Run(cmd, []string{"hel1"})

	expOut := `ID:		3
Name:		hel1
Description:	Helsinki DC Park 1
Network Zone:	eu-central
Country:	FI
City:		Helsinki
Latitude:	60.169855
Longitude:	24.938379
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
