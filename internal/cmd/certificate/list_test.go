package certificate_test

import (
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := certificate.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.CertificateClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.CertificateListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Certificate{
			{
				ID:            123,
				Name:          "test",
				Type:          hcloud.CertificateTypeManaged,
				DomainNames:   []string{"example.com"},
				NotValidAfter: time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
				Created:       time.Now().Add(-20 * time.Minute),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	require.NoError(t, err)
	assert.Empty(t, errOut)
	snaps.MatchSnapshot(t, out)
}
