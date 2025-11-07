package rrset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone/rrset"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := rrset.LabelCmds.AddCobraCommand(fx.State())
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
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
		Return(rrSet, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		UpdateRRSet(gomock.Any(), rrSet, hcloud.ZoneRRSetUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "key=value"})

	expOut := "Label(s) key added to Zone RRSet www A\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := rrset.LabelCmds.RemoveCobraCommand(fx.State())
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
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		GetRRSetByNameAndType(gomock.Any(), z, "www", hcloud.ZoneRRSetTypeA).
		Return(rrSet, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		UpdateRRSet(gomock.Any(), rrSet, hcloud.ZoneRRSetUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"example.com", "www", "A", "key"})

	expOut := "Label(s) key removed from Zone RRSet www A\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
