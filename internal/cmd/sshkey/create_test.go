package sshkey_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := sshkey.CreateCmd.CobraCommand(fx.State())
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

	out, errOut, err := fx.Run(cmd, []string{"--name", "test", "--public-key", "test"})

	expOut := "SSH key 123 created\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := sshkey.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.SSHKeyClient.EXPECT().
		Create(gomock.Any(), hcloud.SSHKeyCreateOpts{
			Name:      "test",
			PublicKey: "test",
			Labels:    make(map[string]string),
		}).
		Return(&hcloud.SSHKey{
			ID:          123,
			Name:        "test",
			PublicKey:   "test",
			Created:     time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
			Labels:      make(map[string]string),
			Fingerprint: "00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00",
		}, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test", "--public-key", "test"})

	expOut := "SSH key 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJSON, jsonOut)
}
