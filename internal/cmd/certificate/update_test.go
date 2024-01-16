package certificate_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateName(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Certificate{ID: 123}, nil, nil)
	fx.Client.CertificateClient.EXPECT().
		Update(gomock.Any(), &hcloud.Certificate{ID: 123}, hcloud.CertificateUpdateOpts{
			Name: "new-name",
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--name", "new-name"})

	expOut := "certificate 123 updated\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
