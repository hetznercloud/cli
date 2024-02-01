package datacenter_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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
  Available:
    No available server types
  Supported:
    No supported server types
`

	assert.NoError(t, err)
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
			ID:          1,
			Name:        "cx11",
			Description: "CX11",
		},
		{
			ID:          3,
			Name:        "cx21",
			Description: "CX21",
		},
		{
			ID:          5,
			Name:        "cx31",
			Description: "CX31",
		},
		{
			ID:          7,
			Name:        "cx41",
			Description: "CX41",
		},
		{
			ID:          9,
			Name:        "cx51",
			Description: "CX51",
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
				Available: serverTypes,
				Supported: serverTypes,
			},
		}, nil, nil)

	for i := 0; i < 2; i++ {
		for _, st := range serverTypes {
			fx.Client.ServerTypeClient.EXPECT().
				ServerTypeName(st.ID).
				Return(st.Name)
			fx.Client.ServerTypeClient.EXPECT().
				ServerTypeDescription(st.ID).
				Return(st.Description)
		}
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
  Available:
  - ID:		 1
    Name:	 cx11
    Description: CX11
  - ID:		 3
    Name:	 cx21
    Description: CX21
  - ID:		 5
    Name:	 cx31
    Description: CX31
  - ID:		 7
    Name:	 cx41
    Description: CX41
  - ID:		 9
    Name:	 cx51
    Description: CX51
  Supported:
  - ID:		 1
    Name:	 cx11
    Description: CX11
  - ID:		 3
    Name:	 cx21
    Description: CX21
  - ID:		 5
    Name:	 cx31
    Description: CX31
  - ID:		 7
    Name:	 cx41
    Description: CX41
  - ID:		 9
    Name:	 cx51
    Description: CX51
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
