package iso

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.ISOClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.ISO{
			ID:           123,
			Name:         "test",
			Description:  "Test ISO",
			Type:         hcloud.ISOTypePublic,
			Architecture: hcloud.Ptr(hcloud.ArchitectureX86),
		}, nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Name:		test
Description:	Test ISO
Type:		public
Architecture:	x86
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
