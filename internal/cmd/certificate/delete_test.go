package certificate_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.DeleteCmd.CobraCommand(fx.State())
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

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "certificate test deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
