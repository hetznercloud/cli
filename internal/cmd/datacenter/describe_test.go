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

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.DatacenterClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Datacenter{
			ID:          4,
			Name:        "fsn1-dc14",
			Location:    &hcloud.Location{Name: "fsn1"},
			Description: "Falkenstein 1 virtual DC 14",
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		4
Name:		fsn1-dc14
Description:	Falkenstein 1 virtual DC 14
Location:
  Name:		fsn1
  Description:	
  Country:	
  City:		
  Latitude:	0.000000
  Longitude:	0.000000
Server Types:
  Available:
    No available server types
  Supported:
    No supported server types
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
