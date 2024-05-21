package servertype_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := servertype.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerType.EXPECT().
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

	out, errOut, err := fx.Run(cmd, []string{"cax11"})

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
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
