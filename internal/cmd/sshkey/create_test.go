package sshkey

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.SSHKeyClient.EXPECT().
		Create(gomock.Any(), hcloud.SSHKeyCreateOpts{
			Name:      "test",
			PublicKey: "test",
			Labels:    make(map[string]string),
		}).
		Return(&hcloud.SSHKey{
			ID:        123,
			Name:      "test",
			PublicKey: "test",
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"--name", "test", "--public-key", "test"})

	expOut := "SSH key 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
