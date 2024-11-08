package certificate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := certificate.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	cert := &hcloud.Certificate{
		ID:             123,
		Name:           "test",
		Type:           hcloud.CertificateTypeManaged,
		Created:        time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
		NotValidBefore: time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
		NotValidAfter:  time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
		DomainNames:    []string{"example.com"},
		Labels:         map[string]string{"key": "value"},
		UsedBy: []hcloud.CertificateUsedByRef{{
			ID:   123,
			Type: hcloud.CertificateUsedByRefTypeLoadBalancer,
		}},
		Status: &hcloud.CertificateStatus{
			Error: &hcloud.Error{
				Code:    hcloud.ErrorCode("cert_error"),
				Message: "Certificate error",
			},
			Issuance: hcloud.CertificateStatusTypeFailed,
			Renewal:  hcloud.CertificateStatusTypeScheduled,
		},
	}

	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(cert, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		LoadBalancerName(int64(123)).
		Return("test")

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:			123
Name:			test
Type:			managed
Fingerprint:		
Created:		%s (%s)
Not valid before:	%s (%s)
Not valid after:	%s (%s)
Status:
  Issuance:	failed
  Renewal:	scheduled
  Failure reason: Certificate error
Domain names:
  - example.com
Labels:
  key:	value
Used By:
  - Type: load_balancer
  - Name: test
`,
		util.Datetime(cert.Created), humanize.Time(cert.Created),
		util.Datetime(cert.NotValidBefore), humanize.Time(cert.NotValidBefore),
		util.Datetime(cert.NotValidAfter), humanize.Time(cert.NotValidAfter))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
