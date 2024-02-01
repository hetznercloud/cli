package servertype_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/cmd/util"
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
				Sort:     []string{"id:asc"},
			},
		).
		Return([]*hcloud.ServerType{
			{
				ID:              123,
				Name:            "test",
				Cores:           2,
				CPUType:         hcloud.CPUTypeShared,
				Architecture:    hcloud.ArchitectureARM,
				Memory:          8.0,
				Disk:            80,
				StorageType:     hcloud.StorageTypeLocal,
				IncludedTraffic: 20 * util.Tebibyte,
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   CORES   CPU TYPE   ARCHITECTURE   MEMORY   DISK    STORAGE TYPE   TRAFFIC
123   test   2       shared     arm            8.0 GB   80 GB   local          20 TB
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
