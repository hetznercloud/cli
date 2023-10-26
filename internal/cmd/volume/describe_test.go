package volume

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Volume{
			ID:     123,
			Name:   "test",
			Size:   50,
			Server: &hcloud.Server{ID: 321},
			Location: &hcloud.Location{
				ID:          3,
				Name:        "hel1",
				Description: "Helsinki DC Park 1",
				NetworkZone: "eu-central",
				Country:     "FI",
				City:        "Helsinki",
				Latitude:    60.169855,
				Longitude:   24.938379,
			},
			Created: time.Date(1905, 10, 6, 12, 0, 0, 0, time.UTC),
		}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Name:		test
Created:	Fri Oct  6 12:00:00 UTC 1905 (a long while ago)
Size:		50 GB
Linux Device:	
Location:
  Name:		hel1
  Description:	Helsinki DC Park 1
  Country:	FI
  City:		Helsinki
  Latitude:	60.169855
  Longitude:	24.938379
Server:
  ID:		321
  Name:		myServer
Protection:
  Delete:	no
Labels:
  No labels
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
