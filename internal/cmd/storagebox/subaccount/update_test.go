package subaccount_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
	}
	sbs := &hcloud.StorageBoxSubaccount{
		ID:         456,
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccountByID(gomock.Any(), sb, int64(456)).
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateSubaccount(gomock.Any(), sbs, hcloud.StorageBoxSubaccountUpdateOpts{
			Description: hcloud.Ptr("new description"),
		})

	out, errOut, err := fx.Run(cmd, []string{"my-storage-box", "456", "--description", "new description"})

	expOut := "Storage Box Subaccount 456 updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
