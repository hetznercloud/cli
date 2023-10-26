package sshkey

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.SSHKey{
			ID:          123,
			Name:        "test",
			Created:     time.Date(1905, 10, 6, 12, 0, 0, 0, time.UTC),
			Fingerprint: "fingerprint",
			PublicKey:   "public key",
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Name:		test
Created:	Fri Oct  6 12:00:00 UTC 1905 (a long while ago)
Fingerprint:	fingerprint
Public Key:
public key
Labels:
  No labels
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
