package sshkey

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sshKey := &hcloud.SSHKey{
		ID:   123,
		Name: "test",
	}

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(sshKey, nil, nil)
	fx.Client.SSHKeyClient.EXPECT().
		Delete(gomock.Any(), sshKey).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "SSH Key test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
