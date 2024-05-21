package iso_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/iso"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := iso.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ISO.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.ISO{
			ID:           123,
			Name:         "test",
			Description:  "Test ISO",
			Type:         hcloud.ISOTypePublic,
			Architecture: hcloud.Ptr(hcloud.ArchitectureX86),
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Name:		test
Description:	Test ISO
Type:		public
Architecture:	x86
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
