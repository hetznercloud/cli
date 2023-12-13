package firewall

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestApplyToServer(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ApplyToResourceCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	firewall := &hcloud.Firewall{
		ID:   123,
		Name: "test",
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(firewall, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(&hcloud.Server{ID: 456}, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		ApplyResources(gomock.Any(), firewall, []hcloud.FirewallResource{{
			Type: hcloud.FirewallResourceTypeServer,
			Server: &hcloud.FirewallResourceServer{
				ID: 456,
			},
		}}).
		Return([]*hcloud.Action{{ID: 123}, {ID: 321}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 321}}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"--type", "server", "--server", "my-server", "test"})

	expOut := "Firewall 123 applied\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestApplyToLabelSelector(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ApplyToResourceCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	firewall := &hcloud.Firewall{
		ID:   123,
		Name: "test",
	}

	fx.Client.FirewallClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(firewall, nil, nil)
	fx.Client.FirewallClient.EXPECT().
		ApplyResources(gomock.Any(), firewall, []hcloud.FirewallResource{{
			Type: hcloud.FirewallResourceTypeLabelSelector,
			LabelSelector: &hcloud.FirewallResourceLabelSelector{
				Selector: "my-label",
			},
		}}).
		Return([]*hcloud.Action{{ID: 123}, {ID: 321}}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 321}}).
		Return(nil)

	out, _, err := fx.Run(cmd, []string{"--type", "label_selector", "--label-selector", "my-label", "test"})

	expOut := "Firewall 123 applied\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
