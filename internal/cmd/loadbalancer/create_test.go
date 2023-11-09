package loadbalancer

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerClient.EXPECT().
		Create(gomock.Any(), hcloud.LoadBalancerCreateOpts{
			Name:             "myLoadBalancer",
			LoadBalancerType: &hcloud.LoadBalancerType{Name: "lb11"},
			Location:         &hcloud.Location{Name: "fsn1"},
			Labels:           make(map[string]string),
		}).
		Return(hcloud.LoadBalancerCreateResult{
			LoadBalancer: &hcloud.LoadBalancer{ID: 123},
			Action:       &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 321}).Return(nil)
	fx.Client.LoadBalancerClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.LoadBalancer{
			ID: 123,
			PublicNet: hcloud.LoadBalancerPublicNet{
				IPv4: hcloud.LoadBalancerPublicNetIPv4{
					IP: net.ParseIP("192.168.2.1"),
				},
				IPv6: hcloud.LoadBalancerPublicNetIPv6{
					IP: net.IPv6zero,
				},
			},
		}, nil, nil)

	out, _, err := fx.Run(cmd, []string{"--name", "myLoadBalancer", "--type", "lb11", "--location", "fsn1"})

	expOut := `Load Balancer 123 created
IPv4: 192.168.2.1
IPv6: ::
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	loadBalancer := &hcloud.LoadBalancer{
		ID: 123,
		PublicNet: hcloud.LoadBalancerPublicNet{
			IPv4: hcloud.LoadBalancerPublicNetIPv4{
				IP: net.ParseIP("192.168.2.1"),
			},
			IPv6: hcloud.LoadBalancerPublicNetIPv6{
				IP: net.IPv6zero,
			},
		},
	}

	fx.Client.LoadBalancerClient.EXPECT().
		Create(gomock.Any(), hcloud.LoadBalancerCreateOpts{
			Name:             "myLoadBalancer",
			LoadBalancerType: &hcloud.LoadBalancerType{Name: "lb11"},
			Location:         &hcloud.Location{Name: "fsn1"},
			Labels:           make(map[string]string),
		}).
		Return(hcloud.LoadBalancerCreateResult{
			LoadBalancer: &hcloud.LoadBalancer{ID: 123},
			Action:       &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 321}).Return(nil)
	fx.Client.LoadBalancerClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(loadBalancer, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		ChangeProtection(gomock.Any(), loadBalancer, hcloud.LoadBalancerChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 333}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 333}).Return(nil)

	out, _, err := fx.Run(cmd, []string{"--name", "myLoadBalancer", "--type", "lb11", "--location", "fsn1", "--enable-protection", "delete"})

	expOut := `Load Balancer 123 created
Resource protection enabled for Load Balancer 123
IPv4: 192.168.2.1
IPv6: ::
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
