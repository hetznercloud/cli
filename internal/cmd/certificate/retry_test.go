package certificate_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRetry(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := certificate.RetryCmd.CobraCommand(fx.State())

	cert := &hcloud.Certificate{
		ID:   123,
		Name: "my-test-cert",
	}

	fx.ExpectEnsureToken()
	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(cert, nil, nil)
	fx.Client.CertificateClient.EXPECT().
		RetryIssuance(gomock.Any(), cert).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123"})

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Retried issuance of certificate my-test-cert\n", out)
}
