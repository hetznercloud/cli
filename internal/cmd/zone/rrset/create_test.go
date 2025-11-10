package rrset_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone/rrset"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := rrset.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		CreateRRSet(gomock.Any(), z, hcloud.ZoneRRSetCreateOpts{
			Name: "www",
			Type: hcloud.ZoneRRSetTypeA,
			Records: []hcloud.ZoneRRSetRecord{{
				Value: "198.51.100.1",
			}},
			Labels: map[string]string{
				"foo": "bar",
			},
			TTL: hcloud.Ptr(42),
		}).
		Return(hcloud.ZoneRRSetCreateResult{
			RRSet: &hcloud.ZoneRRSet{
				ID:   "www/A",
				Type: hcloud.ZoneRRSetTypeA,
				Name: "www",
			},
			Action: &hcloud.Action{ID: 123},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"example.com", "--type", "A", "--name", "www", "--record", "198.51.100.1", "--ttl", "42", "--label", "foo=bar"})

	expOut := "Zone RRSet www A created\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
