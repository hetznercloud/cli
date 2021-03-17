package firewall

import (
	"fmt"
	"net"
)

func validateFirewallIP(ip string) (*net.IPNet, error) {
	i, n, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	if i.String() != n.IP.String() {
		return nil, fmt.Errorf("%s is not the start of the cidr block %s", ip, n)
	}

	return n, nil
}
