package sshkey_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := sshkey.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.SSHKeyClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.SSHKeyListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.SSHKey{
			{
				ID:      123,
				Name:    "test",
				Created: time.Now().Add(-1 * time.Hour),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   FINGERPRINT   AGE
123   test   -             1h 
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
