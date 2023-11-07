package servertype

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cax11").
		Return(&hcloud.ServerType{
			ID:          45,
			Name:        "cax11",
			Description: "CAX11",
			Cores:       2,
			CPUType:     hcloud.CPUTypeShared,
			Memory:      4.0,
			Disk:        40,
			StorageType: hcloud.StorageTypeLocal,
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"cax11"})

	expOut := `ID:			45
Name:			cax11
Description:		CAX11
Cores:			2
CPU Type:		shared
Architecture:		
Memory:			4.0 GB
Disk:			40 GB
Storage Type:		local
Included Traffic:	0 TB
Pricings per Location:
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
