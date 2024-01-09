package sshkey_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/sshkey"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := sshkey.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	key := &hcloud.SSHKey{
		ID:          123,
		Name:        "test",
		Created:     time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
		Fingerprint: "fingerprint",
		PublicKey:   "public key",
	}

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(key, nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Name:		test
Created:	%s (%s)
Fingerprint:	fingerprint
Public Key:
public key
Labels:
  No labels
`, util.Datetime(key.Created), humanize.Time(key.Created))

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
