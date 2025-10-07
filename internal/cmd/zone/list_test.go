package zone_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := zone.ListCmd.CobraCommand(fx.State())

	z := &hcloud.Zone{
		Created: time.Date(2016, time.January, 30, 23, 55, 0, 0, time.UTC),
		ID:      42,
		Labels: map[string]string{
			"environment":    "prod",
			"example.com/my": "label",
			"just-a-key":     "",
		},
		Mode: hcloud.ZoneModePrimary,
		Name: "example.com",
		Protection: hcloud.ZoneProtection{
			Delete: false,
		},
		RecordCount: 0,
		Registrar:   hcloud.ZoneRegistrarHetzner,
		Status:      hcloud.ZoneStatusOk,
		TTL:         10800,
	}

	fx.ExpectEnsureToken()
	fx.Client.ZoneClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ZoneListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Zone{z}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := fmt.Sprintf(`ID   NAME          STATUS   MODE      RECORD COUNT   AGE  
42   example.com   ok       primary   0              %s
`, util.Age(z.Created, time.Now()))

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}

func TestListPrimaryNameservers(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := zone.ListCmd.CobraCommand(fx.State())

	z := &hcloud.Zone{
		Created: time.Date(2024, time.January, 30, 23, 55, 0, 0, time.UTC),
		ID:      42,
		Mode:    hcloud.ZoneModeSecondary,
		Name:    "secondary.example.com",
		PrimaryNameservers: []hcloud.ZonePrimaryNameserver{
			{Address: "primary.example.com", Port: 53},
			{Address: "192.0.2.7", Port: 53},
			{Address: "2001:db8::7", Port: 53},
		},
		TTL: 10800,
	}

	fx.ExpectEnsureToken()
	fx.Client.ZoneClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ZoneListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Zone{z}, nil)

	out, errOut, err := fx.Run(cmd, []string{"-o=columns=id,primary_nameservers"})

	expOut := `ID   PRIMARY NAMESERVERS                                   
42   primary.example.com:53, 192.0.2.7:53, [2001:db8::7]:53
`

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
func TestListAuthoritativeNameservers(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := zone.ListCmd.CobraCommand(fx.State())

	z := &hcloud.Zone{
		ID:   42,
		Mode: hcloud.ZoneModePrimary,
		Name: "primary.example.com",
		AuthoritativeNameservers: hcloud.ZoneAuthoritativeNameservers{
			Assigned: []string{"helium.ns.hetzner.de.", "hydrogen.ns.hetzner.com.", "oxygen.ns.hetzner.com."},
		},
	}

	fx.ExpectEnsureToken()
	fx.Client.ZoneClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ZoneListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.Zone{z}, nil)

	out, errOut, err := fx.Run(cmd, []string{"-o=columns=id,authoritative_nameservers"})

	expOut := `ID   AUTHORITATIVE NAMESERVERS                                              
42   helium.ns.hetzner.de., hydrogen.ns.hetzner.com., oxygen.ns.hetzner.com.
`

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
