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

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storagebox.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.StorageBox{ID: 123}, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		Update(gomock.Any(), &hcloud.StorageBox{ID: 123}, hcloud.StorageBoxUpdateOpts{
			Name: "new-name",
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "Storage Box 123 updated\n"

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
