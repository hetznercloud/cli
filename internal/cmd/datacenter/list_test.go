package datacenter

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

	out, err := fx.Run(cmd, []string{})

	expOut := `ID   NAME        DESCRIPTION                   LOCATION
4    fsn1-dc14   Falkenstein 1 virtual DC 14   fsn1
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
