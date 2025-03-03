package servertype_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := servertype.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.ServerTypeClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ServerTypeListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     nil, // Server Types do not support sorting
			},
		).
		Return([]*hcloud.ServerType{
			{
				ID:           123,
				Name:         "test",
				Cores:        2,
				CPUType:      hcloud.CPUTypeShared,
				Architecture: hcloud.ArchitectureARM,
				Memory:       8.0,
				Disk:         80,
				StorageType:  hcloud.StorageTypeLocal,
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   CORES   CPU TYPE   ARCHITECTURE   MEMORY   DISK    STORAGE TYPE
123   test   2       shared     arm            8.0 GB   80 GB   local       
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestListColumnDeprecated(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := servertype.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.ServerTypeClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ServerTypeListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     nil, // Server Types do not support sorting
			},
		).
		Return([]*hcloud.ServerType{
			{
				ID:   123,
				Name: "deprecated",
				DeprecatableResource: hcloud.DeprecatableResource{
					Deprecation: &hcloud.DeprecationInfo{
						Announced:        time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
						UnavailableAfter: time.Date(2037, 8, 20, 12, 0, 0, 0, time.UTC),
					},
				},
			},
			{
				ID:   124,
				Name: "current",
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"-o=columns=id,name,deprecated"})

	expOut := `ID    NAME         DEPRECATED                  
123   deprecated   Thu Aug 20 12:00:00 UTC 2037
124   current      -                           
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
