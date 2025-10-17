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

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Delete(gomock.Any(), sb).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Storage Box test deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	storageBoxes := []*hcloud.StorageBox{
		{
			ID:   123,
			Name: "test1",
		},
		{
			ID:   456,
			Name: "test2",
		},
		{
			ID:   789,
			Name: "test3",
		},
	}

	var names []string
	for _, sb := range storageBoxes {
		names = append(names, sb.Name)
		fx.Client.StorageBoxClient.EXPECT().
			Get(gomock.Any(), sb.Name).
			Return(sb, nil, nil)
		fx.Client.StorageBoxClient.EXPECT().
			Delete(gomock.Any(), sb).
			Return(&hcloud.Action{ID: sb.ID + 1000}, nil, nil)
	}

	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 1123}, &hcloud.Action{ID: 1456}, &hcloud.Action{ID: 1789}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Storage Boxes test1, test2, test3 deleted\n", out)
}
