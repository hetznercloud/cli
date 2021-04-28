package network

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := hcapi2.NewMockClient(ctrl)
	tokenEnsurer := state.NewMockTokenEnsurer(ctrl)

	cmd := ListCmd.CobraCommand(context.Background(), client, tokenEnsurer)

	tokenEnsurer.EXPECT().EnsureToken(gomock.Any(), gomock.Any()).Return(nil)
	client.NetworkClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.NetworkListOpts{
				ListOpts: hcloud.ListOpts{
					PerPage:       50,
					LabelSelector: "foo=bar",
				},
			},
		).
		Return([]*hcloud.Network{
			{
				ID:      123,
				Name:    "test-net",
				IPRange: &net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: net.CIDRMask(24, 32)},
				Servers: []*hcloud.Server{{ID: 3421}},
			},
		},
			nil)

	args := []string{"--selector", "foo=bar"}
	cmd.SetArgs(args)

	out, err := testutil.CaptureStdout(cmd.Execute)

	expOut := `ID    NAME       IP RANGE       SERVERS
123   test-net   192.0.2.1/24   1 server
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
