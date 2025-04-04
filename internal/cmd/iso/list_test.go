package iso_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/iso"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := iso.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.ISOClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ISOListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     nil, // ISOs do not support sorting
			},
		).
		Return([]*hcloud.ISO{
			{
				ID:           123,
				Name:         "test",
				Description:  "Test ISO",
				Type:         hcloud.ISOTypePublic,
				Architecture: hcloud.Ptr(hcloud.ArchitectureX86),
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   DESCRIPTION   TYPE     ARCHITECTURE
123   test   Test ISO      public   x86         
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
