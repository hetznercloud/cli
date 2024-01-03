package primaryip

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAssign(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := AssignCmd.CobraCommand(fx.State())
	action := &hcloud.Action{ID: 1}
	fx.ExpectEnsureToken()
	var (
		server = hcloud.Server{
			ID:   15,
			Name: "my server",
		}
	)

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), fmt.Sprintf("%d", server.ID)).
		Return(&server, nil, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			&hcloud.PrimaryIP{ID: 13},
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		Assign(
			gomock.Any(),
			hcloud.PrimaryIPAssignOpts{
				ID:           13,
				AssigneeType: "server",
				AssigneeID:   15,
			},
		).
		Return(
			action,
			&hcloud.Response{},
			nil,
		)

	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), gomock.Any(), action).Return(nil)

	out, _, err := fx.Run(cmd, []string{"13", "--server", "15"})

	expOut := "Primary IP 13 assigned to server 15\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
