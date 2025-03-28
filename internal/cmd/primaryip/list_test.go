package primaryip_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.PrimaryIPListOpts{
				ListOpts: hcloud.ListOpts{
					PerPage:       50,
					LabelSelector: "foo=bar",
				},
				Sort: []string{"id:asc"},
			},
		).
		Return([]*hcloud.PrimaryIP{
			{
				ID:         123,
				Name:       "test-net",
				AutoDelete: true,
				Type:       hcloud.PrimaryIPTypeIPv4,
				IP:         net.ParseIP("127.0.0.1"),
				Created:    time.Now().Add(-10 * time.Second),
			},
		},
			nil)

	out, errOut, err := fx.Run(cmd, []string{"--selector", "foo=bar"})

	expOut := `ID    TYPE   NAME       IP          ASSIGNEE   DNS   AUTO DELETE   AGE
123   ipv4   test-net   127.0.0.1   -          -     yes           10s
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
