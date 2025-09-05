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

func TestFolders(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.FoldersCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.StorageBox{ID: 123}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Folders(gomock.Any(), &hcloud.StorageBox{ID: 123}, hcloud.StorageBoxFoldersOpts{Path: ""}).
		Return(hcloud.StorageBoxFoldersResult{
			Folders: []string{"folder1", "folder2"},
		}, nil, nil)

	args := []string{"test"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, `Folders:
- folder1
- folder2
`, out)
}

func TestFoldersJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.FoldersCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.StorageBox{ID: 123}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Folders(gomock.Any(), &hcloud.StorageBox{ID: 123}, hcloud.StorageBoxFoldersOpts{Path: ""}).
		Return(hcloud.StorageBoxFoldersResult{
			Folders: []string{"folder1", "folder2"},
		}, nil, nil)

	args := []string{"test", "-o=json"}
	jsonOut, out, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, out)
	assert.JSONEq(t, `{"folders": ["folder1", "folder2"]}`, jsonOut)
}

func TestFoldersYAML(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.FoldersCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.StorageBox{ID: 123}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Folders(gomock.Any(), &hcloud.StorageBox{ID: 123}, hcloud.StorageBoxFoldersOpts{Path: ""}).
		Return(hcloud.StorageBoxFoldersResult{
			Folders: []string{"folder1", "folder2"},
		}, nil, nil)

	args := []string{"test", "-o=yaml"}
	yamlOut, out, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, out)
	assert.YAMLEq(t, `folders:
- folder1
- folder2
`, yamlOut)
}
