package certificate

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Certificate{
			ID:             123,
			Name:           "test",
			Type:           hcloud.CertificateTypeManaged,
			Created:        time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
			NotValidBefore: time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
			NotValidAfter:  time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
			DomainNames:    []string{"example.com"},
			Labels:         map[string]string{"key": "value", "key2": "value2"},
			UsedBy: []hcloud.CertificateUsedByRef{{
				ID:   123,
				Type: hcloud.CertificateUsedByRefTypeLoadBalancer,
			}},
		}, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		LoadBalancerName(int64(123)).
		Return("test")

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:			123
Name:			test
Type:			managed
Fingerprint:		
Created:		Mon Aug 24 12:00:00 UTC 2020 (3 years ago)
Not valid before:	Mon Aug 24 12:00:00 UTC 2020 (3 years ago)
Not valid after:	Wed Aug 20 12:00:00 UTC 2036 (12 years from now)
Domain names:
  - example.com
Labels:
  key:	value
  key2:	value2
Used By:
  - Type: load_balancer
  - Name: test
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
