package sshkey_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.SSHKey{ID: 123}, nil, nil)
	fx.Client.SSHKeyClient.EXPECT().
		Update(gomock.Any(), &hcloud.SSHKey{ID: 123}, hcloud.SSHKeyUpdateOpts{
			Name: "new-name",
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "SSHKey 123 updated\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
