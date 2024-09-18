//go:build e2e

package e2e

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCombined(t *testing.T) {
	// combined tests combine multiple resources and can thus not be run in parallel
	serverName := withSuffix("test-server")
	serverID, err := createServer(t, serverName, TestServerType, TestImage)
	if err != nil {
		t.Fatal(err)
	}

	firewallName := withSuffix("test-firewall")
	firewallID, err := createFirewall(t, firewallName)
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", serverName, strconv.FormatInt(firewallID, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d applied to resource\n", firewallID), out)

	out, err = runCommand(t, "firewall", "delete", strconv.FormatInt(firewallID, 10))
	assert.Regexp(t, `^firewall with ID [0-9]+ is still in use \(resource_in_use, [0-9a-f]+\)$`, err.Error())
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", serverName, strconv.FormatInt(firewallID, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d removed from resource\n", firewallID), out)

	out, err = runCommand(t, "firewall", "delete", strconv.FormatInt(firewallID, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallID), out)

	floatingIPName := withSuffix("test-floating-ip")
	floatingIP, err := createFloatingIP(t, floatingIPName, "ipv4", "--server", strconv.FormatInt(serverID, 10))
	if err != nil {
		t.Fatal(err)
	}

	out, err = runCommand(t, "floating-ip", "unassign", strconv.FormatInt(floatingIP, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d unassigned\n", floatingIP), out)

	out, err = runCommand(t, "floating-ip", "assign", strconv.FormatInt(floatingIP, 10), strconv.FormatInt(serverID, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d assigned to server %d\n", floatingIP, serverID), out)

	out, err = runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIP, 10))
	require.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv4
Name:\s+test-floating-ip-[0-9a-f]{8}
Description:\s+-
Created:.*?
IP:\s+(?:[0-9]{1,3}\.){3}[0-9]{1,3}
Blocked:\s+no
Home Location:\s+[a-z]{3}[0-9]*
Server:
\s+ID:\s+[0-9]+
\s+Name:\s+test-server-[0-9a-f]{8}
DNS:
.*?
Protection:
\s+Delete:\s+no
Labels:
\s+No labels
`, out)

	out, err = runCommand(t, "floating-ip", "list", "-o", "columns=server", "-o", "noheader")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s\n", serverName), out)

	out, err = runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIP, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIP), out)

	out, err = runCommand(t, "server", "delete", strconv.FormatInt(serverID, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Server %d deleted\n", serverID), out)
}
