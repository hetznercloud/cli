package firewall

import (
	"fmt"
	"net"
	"reflect"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func ValidateFirewallIP(ip string) (*net.IPNet, error) {
	i, n, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, err
	}
	if i.String() != n.IP.String() {
		return nil, fmt.Errorf("%s is not the start of the cidr block %s", ip, n)
	}

	return n, nil
}

func EqualFirewallRule(a, b hcloud.FirewallRule) bool {
	a.Description = nil
	b.Description = nil

	return reflect.DeepEqual(a, b)
}
