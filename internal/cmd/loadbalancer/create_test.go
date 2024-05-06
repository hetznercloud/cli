package loadbalancer_test

import (
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.CreateCmd.CobraCommand(fx.State())
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
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321}).Return(nil)
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

	out, errOut, err := fx.Run(cmd, []string{"--name", "myLoadBalancer", "--type", "lb11", "--location", "fsn1"})

	expOut := `Load Balancer 123 created
IPv4: 192.168.2.1
IPv6: ::
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	lb := &hcloud.LoadBalancer{
		ID:   123,
		Name: "myLoadBalancer",
		PublicNet: hcloud.LoadBalancerPublicNet{
			IPv4: hcloud.LoadBalancerPublicNetIPv4{
				IP: net.ParseIP("192.168.2.1"),
			},
			IPv6: hcloud.LoadBalancerPublicNetIPv6{
				IP: net.IPv6zero,
			},
		},
		Labels:          make(map[string]string),
		Created:         time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
		IncludedTraffic: 654321,
		Services:        []hcloud.LoadBalancerService{},
		Targets:         []hcloud.LoadBalancerTarget{},
	}

	fx.Client.LoadBalancerClient.EXPECT().
		Create(gomock.Any(), hcloud.LoadBalancerCreateOpts{
			Name:             "myLoadBalancer",
			LoadBalancerType: &hcloud.LoadBalancerType{Name: "lb11"},
			Location:         &hcloud.Location{Name: "fsn1"},
			Labels:           make(map[string]string),
		}).
		Return(hcloud.LoadBalancerCreateResult{
			LoadBalancer: lb,
			Action:       &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321}).Return(nil)
	fx.Client.LoadBalancerClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(lb, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "myLoadBalancer", "--type", "lb11", "--location", "fsn1"})

	expOut := "Load Balancer 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}

func TestCreateProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.CreateCmd.CobraCommand(fx.State())
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
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321}).Return(nil)
	fx.Client.LoadBalancerClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(loadBalancer, nil, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		ChangeProtection(gomock.Any(), loadBalancer, hcloud.LoadBalancerChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 333}, nil, nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 333}).Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "myLoadBalancer", "--type", "lb11", "--location", "fsn1", "--enable-protection", "delete"})

	expOut := `Load Balancer 123 created
Resource protection enabled for Load Balancer 123
IPv4: 192.168.2.1
IPv6: ::
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
