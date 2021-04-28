package location

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestDescribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := hcapi2.NewMockClient(ctrl)
	actionWaiter := state.NewMockActionWaiter(ctrl)
	tokenEnsurer := state.NewMockTokenEnsurer(ctrl)

	cmd := newDescribeCommand(context.Background(), client, tokenEnsurer, actionWaiter)

	tokenEnsurer.EXPECT().EnsureToken(gomock.Any(), gomock.Any()).Return(nil)
	client.LocationClient.EXPECT().
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

	args := []string{"hel1"}
	cmd.SetArgs(args)

	out, err := testutil.CaptureStdout(func() error {
		return cmd.Execute()
	})

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
