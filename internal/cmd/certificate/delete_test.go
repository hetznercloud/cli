package certificate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

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

	expOut := "Certificate test deleted\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	certs := []*hcloud.Certificate{
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

	var names []string
	for _, cert := range certs {
		names = append(names, cert.Name)
		fx.Client.CertificateClient.EXPECT().
			Get(gomock.Any(), cert.Name).
			Return(cert, nil, nil)
		fx.Client.CertificateClient.EXPECT().
			Delete(gomock.Any(), cert).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Certificates test1, test2, test3 deleted\n", out)
}
