package rrset_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone/rrset"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := rrset.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
	}

	rrSet := &hcloud.ZoneRRSet{
		Zone: z,
		ID:   "www/A",
		Name: "www",
		Type: hcloud.ZoneRRSetTypeA,
		TTL:  hcloud.Ptr(600),
		Labels: map[string]string{
			"environment":    "prod",
			"example.com/my": "label",
			"just-a-key":     "",
		},
		Records: []hcloud.ZoneRRSetRecord{
			{Value: "198.51.100.1", Comment: "My web server at Hetzner Cloud."},
			{Value: "198.51.100.2"},
		},
		Protection: hcloud.ZoneRRSetProtection{
			Change: true,
		},
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
		Return(rrSet, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A"})

	expOut := `ID:    www/A
Type:  A
Name:  www
TTL:   600

Protection:
  Change:  yes

Labels:
  environment:     prod
  example.com/my:  label
  just-a-key:      

Records:
  - Value:    198.51.100.1
    Comment:  My web server at Hetzner Cloud.
  - Value:    198.51.100.2
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
