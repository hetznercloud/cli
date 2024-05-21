package volume_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := volume.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.Volume.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.VolumeListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Volume{
			{
				ID:       123,
				Name:     "test",
				Size:     50,
				Server:   &hcloud.Server{ID: 321},
				Location: &hcloud.Location{Name: "fsn1"},
				Created:  time.Now().Add(-1 * time.Hour),
			},
		}, nil)
	fx.Client.Server.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   SIZE    SERVER     LOCATION   AGE
123   test   50 GB   myServer   fsn1       1h
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
