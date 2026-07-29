package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/deprecationutil"
)

const ChangeDeprecatedLoadBalancerTypeMessage = `Existing Load Balancers of that plan will ` +
	`continue to work as before and no action is required on your part. ` +
	`It is possible to migrate this Load Balancer to another Load Balancer Type by using ` +
	`the "hcloud load-balancer change-type" command.`

func deprecatedLoadBalancerTypeWarning(loadBalancerType *hcloud.LoadBalancerType) string {
	message, _ := deprecationutil.LoadBalancerTypeMessage(loadBalancerType)
	if message == "" {
		return ""
	}

	return fmt.Sprintf("Attention: %s. %s\n\n", message, ChangeDeprecatedLoadBalancerTypeMessage)
}
