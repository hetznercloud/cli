package iso

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := ListCmd.CobraCommand(context.Background(), fx.Client, fx.TokenEnsurer)

	fx.ExpectEnsureToken()
	fx.Client.ISOClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ISOListOpts{
				ListOpts: hcloud.ListOpts{PerPage: 50},
				Sort:     []string{"id:asc"},
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

	out, _, err := fx.Run(cmd, []string{})

	expOut := `ID    NAME   DESCRIPTION   TYPE     ARCHITECTURE
123   test   Test ISO      public   x86
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
