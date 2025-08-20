package storageboxtype_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storageboxtype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storageboxtype.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxTypeClient.EXPECT().
		Get(gomock.Any(), "bx11").
		Return(&hcloud.StorageBoxType{
			ID:                     42,
			Name:                   "bx11",
			Description:            "BX11",
			SnapshotLimit:          10,
			AutomaticSnapshotLimit: 10,
			SubaccountsLimit:       200,
			Size:                   1073741824,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"bx11"})

	expOut := `ID:			42
Name:			bx11
Description:		BX11
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
