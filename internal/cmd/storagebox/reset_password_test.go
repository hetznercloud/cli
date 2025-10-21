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

func TestResetPassword(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.ResetPasswordCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{ID: 123, Name: "my-storage-box"}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		ResetPassword(gomock.Any(), sb, hcloud.StorageBoxResetPasswordOpts{Password: "new-password"}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	args := []string{"my-storage-box", "--password", "new-password"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, "Password of Storage Box 123 reset\n", out)
}
