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

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := rrset.ListCmd.CobraCommand(fx.State())

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

	fx.ExpectEnsureToken()
	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "example.com").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		AllRRSetsWithOpts(
			gomock.Any(),
			z,
			hcloud.ZoneRRSetListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.ZoneRRSet{rrSet}, nil)

	out, errOut, err := fx.Run(cmd, []string{"example.com", "-o=columns=name,type,ttl,protection,records,labels"})

	expOut := `NAME   TYPE   TTL   PROTECTION   RECORDS        LABELS                                            
www    A      600   change       198.51.100.1   environment=prod, example.com/my=label, just-a-key
                                 198.51.100.2                                                     
`

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
