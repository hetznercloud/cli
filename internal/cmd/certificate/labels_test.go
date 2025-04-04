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

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Certificate{ID: 123}, nil, nil)
	fx.Client.CertificateClient.EXPECT().
		Update(gomock.Any(), &hcloud.Certificate{ID: 123}, hcloud.CertificateUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Certificate 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	cert := &hcloud.Certificate{
		ID: 123,
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(cert, nil, nil)
	fx.Client.CertificateClient.EXPECT().
		Update(gomock.Any(), cert, hcloud.CertificateUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Certificate 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
