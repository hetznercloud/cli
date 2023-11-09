package certificate

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DeleteCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	cert := &hcloud.Certificate{
		ID:   123,
		Name: "test",
	}

	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(cert, nil, nil)
	fx.Client.CertificateClient.EXPECT().
		Delete(gomock.Any(), cert).
		Return(nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := "certificate test deleted\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
