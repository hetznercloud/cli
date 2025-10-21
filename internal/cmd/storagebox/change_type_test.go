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

func TestChangeType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.ChangeTypeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{ID: 123, Name: "my-storage-box"}
	sbt := &hcloud.StorageBoxType{ID: 456, Name: "bx21"}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxTypeClient.EXPECT().
		Get(gomock.Any(), "bx21").
		Return(sbt, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		ChangeType(gomock.Any(), sb, hcloud.StorageBoxChangeTypeOpts{
			StorageBoxType: sbt,
		}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-storage-box", "bx21"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, "Storage Box 123 upgraded to type bx21\n", out)
}
