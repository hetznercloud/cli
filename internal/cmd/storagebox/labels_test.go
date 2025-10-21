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

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Update(gomock.Any(), sb, hcloud.StorageBoxUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Storage Box test\n"

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Update(gomock.Any(), sb, hcloud.StorageBoxUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Storage Box test\n"

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
