package loadbalancer

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerHealth(t *testing.T) {
	tests := []struct {
		name     string
		lb       *hcloud.LoadBalancer
		expected string
	}{
		{
			name: "healthy",
			lb: &hcloud.LoadBalancer{
				Name:     "foobar",
				Services: make([]hcloud.LoadBalancerService, 1),
				Targets: []hcloud.LoadBalancerTarget{
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
						},
					},
				},
			},
			expected: string(hcloud.LoadBalancerTargetHealthStatusStatusHealthy),
		},
		{
			name: "unhealthy",
			lb: &hcloud.LoadBalancer{
				Name:     "foobar",
				Services: make([]hcloud.LoadBalancerService, 1),
				Targets: []hcloud.LoadBalancerTarget{
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
						},
					},
				},
			},
			expected: string(hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy),
		},
		{
			name: "unknown",
			lb: &hcloud.LoadBalancer{
				Name:     "foobar",
				Services: make([]hcloud.LoadBalancerService, 1),
				Targets: []hcloud.LoadBalancerTarget{
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
						},
					},
				},
			},
			expected: string(hcloud.LoadBalancerTargetHealthStatusStatusUnknown),
		},
		{
			name: "mixed",
			lb: &hcloud.LoadBalancer{
				Name:     "foobar",
				Services: make([]hcloud.LoadBalancerService, 1),
				Targets: []hcloud.LoadBalancerTarget{
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
						},
					},
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
						},
					},
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
						},
					},
				},
			},
			expected: "mixed",
		},
		{
			name: "mixed_many_services_grouped_by_target",
			lb: &hcloud.LoadBalancer{
				Name:     "foobar",
				Services: make([]hcloud.LoadBalancerService, 3),
				Targets: []hcloud.LoadBalancerTarget{
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
						},
					},
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
						},
					},
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
						},
					},
				},
			},
			expected: "mixed",
		},
		{
			name: "mixed_many_services_mixed",
			lb: &hcloud.LoadBalancer{
				Name:     "foobar",
				Services: make([]hcloud.LoadBalancerService, 3),
				Targets: []hcloud.LoadBalancerTarget{
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
						},
					},
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
						},
					},
					{
						HealthStatus: []hcloud.LoadBalancerTargetHealthStatus{
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusHealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnhealthy,
							},
							{
								Status: hcloud.LoadBalancerTargetHealthStatusStatusUnknown,
							},
						},
					},
				},
			},
			expected: "mixed",
		},
	}

	for _, test := range tests {
		res := loadBalancerHealth(test.lb)
		assert.Equal(t, test.expected, res, test.name)
	}
}
