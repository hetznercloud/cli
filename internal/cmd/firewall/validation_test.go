package firewall_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/firewall"
)

func TestValidateFirewallIP(t *testing.T) {
	tests := []struct {
		name   string
		ip     string
		errMsg string
	}{
		{
			name: "Valid CIDR (IPv4)",
			ip:   "10.0.0.0/8",
		},
		{
			name: "Valid CIDR (IPv6)",
			ip:   "fe80::/128",
		},
		{
			name:   "Invalid IP",
			ip:     "test",
			errMsg: "invalid CIDR address: test",
		},
		{
			name:   "Missing CIDR notation (IPv4)",
			ip:     "10.0.0.0",
			errMsg: "invalid CIDR address: 10.0.0.0",
		},
		{
			name: "Missing CIDR notation (IPv6)",
			ip:   "fe80::",
			//nolint:revive
			errMsg: "invalid CIDR address: fe80::",
		},
		{
			name:   "Host bit set (IPv4)",
			ip:     "10.0.0.5/8",
			errMsg: "10.0.0.5/8 is not the start of the cidr block 10.0.0.0/8",
		},
		{
			name:   "Host bit set (IPv6)",
			ip:     "fe80::1337/64",
			errMsg: "fe80::1337/64 is not the start of the cidr block fe80::/64",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			net, err := firewall.ValidateFirewallIP(test.ip)

			if test.errMsg != "" {
				require.EqualError(t, err, test.errMsg)
				assert.Nil(t, net)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, net)
		})
	}
}
