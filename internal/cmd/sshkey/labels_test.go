package sshkey_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.SSHKey{ID: 123}, nil, nil)
	fx.Client.SSHKeyClient.EXPECT().
		Update(gomock.Any(), &hcloud.SSHKey{ID: 123}, hcloud.SSHKeyUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to SSH Key 123\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.SSHKey{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)
	fx.Client.SSHKeyClient.EXPECT().
		Update(gomock.Any(), &hcloud.SSHKey{ID: 123}, hcloud.SSHKeyUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from SSH Key 123\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
