package firewall

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.AddCobraCommand(fx.State())
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

	out, _, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to firewall 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Firewall{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		Update(gomock.Any(), &hcloud.Firewall{ID: 123}, hcloud.FirewallUpdateOpts{
			Labels: make(map[string]string),
		})

	out, _, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from firewall 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
