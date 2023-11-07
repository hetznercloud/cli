package certificate

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreateManaged(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		CreateCertificate(gomock.Any(), hcloud.CertificateCreateOpts{
			Name:        "test",
			Type:        hcloud.CertificateTypeManaged,
			DomainNames: []string{"example.com"},
		}).
		Return(hcloud.CertificateCreateResult{
			Certificate: &hcloud.Certificate{
				ID:          123,
				Name:        "test",
				Type:        hcloud.CertificateTypeManaged,
				DomainNames: []string{"example.com"},
			},
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})

	out, err := fx.Run(cmd, []string{"--name", "test", "--type", "managed", "--domain", "example.com"})

	expOut := "Certificate 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateUploaded(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		Create(gomock.Any(), hcloud.CertificateCreateOpts{
			Name:        "test",
			Type:        hcloud.CertificateTypeUploaded,
			Certificate: "certificate file content",
			PrivateKey:  "key file content",
		}).
		Return(&hcloud.Certificate{
			ID:   123,
			Name: "test",
			Type: hcloud.CertificateTypeUploaded,
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"--name", "test", "--key-file", "testdata/key.pem", "--cert-file", "testdata/cert.pem"})

	expOut := "Certificate 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
