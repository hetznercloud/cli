package volume_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := volume.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	v := &hcloud.Volume{
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
		Return(v, nil, nil)
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, errOut, err := fx.Run(cmd, []string{"test"})

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
`, util.Datetime(v.Created), humanize.Time(v.Created))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
