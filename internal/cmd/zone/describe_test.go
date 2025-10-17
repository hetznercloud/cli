package zone_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := zone.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		AuthoritativeNameservers: hcloud.ZoneAuthoritativeNameservers{
			Assigned: []string{
				"hydrogen.ns.hetzner.com.",
				"oxygen.ns.hetzner.com.",
				"helium.ns.hetzner.de.",
			},
			Delegated: []string{
				"hydrogen.ns.hetzner.com.",
				"oxygen.ns.hetzner.com.",
				"helium.ns.hetzner.de.",
			},
			DelegationLastCheck: time.Date(2016, time.January, 30, 23, 55, 0, 0, time.UTC),
			DelegationStatus:    hcloud.ZoneDelegationStatusValid,
		},
		Created: time.Date(2016, time.January, 30, 23, 55, 0, 0, time.UTC),
		ID:      42,
		Labels: map[string]string{
			"environment":    "prod",
			"example.com/my": "label",
			"just-a-key":     "",
		},
		Mode: hcloud.ZoneModePrimary,
		Name: "example.com",
		PrimaryNameservers: []hcloud.ZonePrimaryNameserver{
			{
				Address: "198.51.100.1",
				Port:    53,
			},
			{
				Address: "203.0.113.1",
				Port:    53,
			},
		},
		Protection: hcloud.ZoneProtection{
			Delete: false,
		},
		RecordCount: 0,
		Registrar:   hcloud.ZoneRegistrarHetzner,
		Status:      hcloud.ZoneStatusOk,
		TTL:         10800,
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"example.com"})

	expOut := fmt.Sprintf(`ID:		42
Name:		example.com
Created:	2016-01-30 23:55:00 UTC (%s)
Mode:		primary
Status:		ok
TTL:		10800
Registrar:	hetzner
Record Count:	0
Protection:
  Delete:	no
Labels:
  environment: prod
  example.com/my: label
  just-a-key: 
Authoritative Nameservers:
  Assigned:
    - hydrogen.ns.hetzner.com.
    - oxygen.ns.hetzner.com.
    - helium.ns.hetzner.de.
  Delegated:
    - hydrogen.ns.hetzner.com.
    - oxygen.ns.hetzner.com.
    - helium.ns.hetzner.de.
  Delegation last check:	2016-01-30 23:55:00 UTC (%s)
  Delegation status:		valid
`, humanize.Time(z.Created), humanize.Time(z.AuthoritativeNameservers.DelegationLastCheck))

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
