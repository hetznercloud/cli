package storagebox_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDisableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.DisableProtectionCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.StorageBox{ID: 123}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		ChangeProtection(gomock.Any(), &hcloud.StorageBox{ID: 123}, hcloud.StorageBoxChangeProtectionOpts{
			Delete: false,
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "delete"})

	expOut := "Resource protection disabled for Storage Box 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
