package datacenter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/datacenter"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := datacenter.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.DatacenterClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Datacenter{
			ID:          4,
			Name:        "fsn1-dc14",
			Location:    &hcloud.Location{Name: "fsn1"},
			Description: "Falkenstein 1 virtual DC 14",
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		4
Name:		fsn1-dc14
Description:	Falkenstein 1 virtual DC 14
Location:
  Name:		fsn1
  Description:	
  Country:	
  City:		
  Latitude:	0.000000
  Longitude:	0.000000
Server Types:
  No server types
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDescribeWithTypes(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := datacenter.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	serverTypes := []*hcloud.ServerType{
		{
			ID:          3,
			Name:        "cx22",
			Description: "CX22",
		},
		{
			ID:          5,
			Name:        "cx32",
			Description: "CX32",
		},
		{
			ID:          7,
			Name:        "cx42",
			Description: "CX42",
		},
		{
			ID:          9,
			Name:        "cx52",
			Description: "CX52",
		},
	}

	fx.Client.DatacenterClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Datacenter{
			ID:          4,
			Name:        "fsn1-dc14",
			Location:    &hcloud.Location{Name: "fsn1"},
			Description: "Falkenstein 1 virtual DC 14",
			ServerTypes: hcloud.DatacenterServerTypes{
				Supported: serverTypes,
				Available: serverTypes[:2],
			},
		}, nil, nil)

	for _, st := range serverTypes {
		fx.Client.ServerTypeClient.EXPECT().
			ServerTypeName(st.ID).
			Return(st.Name)
	}

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		4
Name:		fsn1-dc14
Description:	Falkenstein 1 virtual DC 14
Location:
  Name:		fsn1
  Description:	
  Country:	
  City:		
  Latitude:	0.000000
  Longitude:	0.000000
Server Types:
  - ID: 3        Name: cx22     Supported: true     Available: true
  - ID: 5        Name: cx32     Supported: true     Available: true
  - ID: 7        Name: cx42     Supported: true     Available: false
  - ID: 9        Name: cx52     Supported: true     Available: false
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
