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

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := ListCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer)

	fx.ExpectEnsureToken()
	fx.Client.VolumeClient.EXPECT().
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
	fx.Client.ServerClient.EXPECT().
		ServerName(int64(321)).
		Return("myServer")

	out, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   SIZE    SERVER     LOCATION   AGE
123   test   50 GB   myServer   fsn1       1h
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
