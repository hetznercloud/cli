package firewall_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Firewall{ID: 123}, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		Update(gomock.Any(), &hcloud.Firewall{ID: 123}, hcloud.FirewallUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Firewall 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := firewall.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fw := &hcloud.Firewall{
		ID: 123,
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(fw, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		Update(gomock.Any(), fw, hcloud.FirewallUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Firewall 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
