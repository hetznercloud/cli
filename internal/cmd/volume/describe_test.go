package volume

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
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

	volume := &hcloud.Volume{
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
		Created: time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
	}

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(volume, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Name:		test
Created:	%s (%s)
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
`, util.Datetime(volume.Created), humanize.Time(volume.Created))

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
