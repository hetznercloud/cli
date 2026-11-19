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

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}
	sbs := &hcloud.StorageBoxSubaccount{
		ID:         456,
		Name:       "my-subaccount",
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccount(gomock.Any(), sb, "my-subaccount").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateSubaccount(gomock.Any(), sbs, hcloud.StorageBoxSubaccountUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "my-subaccount", "key=value"})

	expOut := "Label(s) key added to Storage Box Subaccount my-subaccount\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := subaccount.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}
	sbs := &hcloud.StorageBoxSubaccount{
		ID:   456,
		Name: "my-subaccount",
		Labels: map[string]string{
			"key": "value",
		},
		StorageBox: sb,
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccount(gomock.Any(), sb, "my-subaccount").
		Return(sbs, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		UpdateSubaccount(gomock.Any(), sbs, hcloud.StorageBoxSubaccountUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "my-subaccount", "key"})

	expOut := "Label(s) key removed from Storage Box Subaccount my-subaccount\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
