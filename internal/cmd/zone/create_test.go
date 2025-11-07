package zone_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ZoneClient.EXPECT().
		Create(gomock.Any(), hcloud.ZoneCreateOpts{
			Name:   "example.com",
			Mode:   hcloud.ZoneModePrimary,
			TTL:    hcloud.Ptr(600),
			Labels: map[string]string{"foo": "bar"},
		}).
		Return(hcloud.ZoneCreateResult{
			Zone: &hcloud.Zone{
				ID:   123,
				Name: "example.com",
			},
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})
	fx.Client.ZoneClient.EXPECT().GetByID(gomock.Any(), int64(123)).Return(
		&hcloud.Zone{
			ID:   123,
			Name: "example.com",
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "example.com", "--mode", "primary", "--ttl", "600", "--label", "foo=bar"})

	expOut := "Zone example.com created\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	testZone := &hcloud.Zone{
		ID:                 123,
		Name:               "example.com",
		Mode:               hcloud.ZoneModePrimary,
		TTL:                600,
		Labels:             map[string]string{"foo": "bar"},
		PrimaryNameservers: nil,
		Created:            time.Date(2016, time.January, 30, 23, 50, 0, 0, time.UTC),
		Protection:         hcloud.ZoneProtection{Delete: false},
		Status:             hcloud.ZoneStatusOk,
		AuthoritativeNameservers: hcloud.ZoneAuthoritativeNameservers{
			Assigned:            nil,
			Delegated:           nil,
			DelegationLastCheck: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
			DelegationStatus:    hcloud.ZoneDelegationStatusUnknown,
		},
		RecordCount: 0,
		Registrar:   hcloud.ZoneRegistrarUnknown,
	}

	fx.Client.ZoneClient.EXPECT().
		Create(gomock.Any(), hcloud.ZoneCreateOpts{
			Name:   "example.com",
			Mode:   hcloud.ZoneModePrimary,
			TTL:    hcloud.Ptr(600),
			Labels: map[string]string{"foo": "bar"},
		}).
		Return(hcloud.ZoneCreateResult{
			Zone:   testZone,
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}})
	fx.Client.ZoneClient.EXPECT().GetByID(gomock.Any(), int64(123)).Return(
		testZone, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "example.com", "--mode", "primary", "--ttl", "600", "--label", "foo=bar"})

	expOut := "Zone example.com created\n"

	require.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJSON, jsonOut)
}
