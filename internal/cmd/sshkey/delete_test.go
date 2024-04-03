package sshkey_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.DeleteCmd.CobraCommand(fx.State())
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

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "SSH Key test deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	keys := []*hcloud.SSHKey{
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

	expOutBuilder := strings.Builder{}

	var names []string
	for _, key := range keys {
		names = append(names, key.Name)
		expOutBuilder.WriteString(fmt.Sprintf("SSH Key %s deleted\n", key.Name))
		fx.Client.SSHKeyClient.EXPECT().
			Get(gomock.Any(), key.Name).
			Return(key, nil, nil)
		fx.Client.SSHKeyClient.EXPECT().
			Delete(gomock.Any(), key).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
