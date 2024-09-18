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
	serverID := createServer(t, "test-server", TestServerType, TestImage)

	firewallID, err := createFirewall(t, "test-firewall")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", "test-server", strconv.Itoa(firewallID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d applied to resource\n", firewallID), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallID))
	assert.Regexp(t, `^firewall with ID [0-9]+ is still in use \(resource_in_use, [0-9a-f]+\)$`, err.Error())
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "test-server", strconv.Itoa(firewallID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d removed from resource\n", firewallID), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallID), out)

	floatingIP, err := createFloatingIP(t, "ipv4", "--server", strconv.Itoa(serverID))
	if err != nil {
		t.Fatal(err)
	}

	out, err = runCommand(t, "floating-ip", "unassign", strconv.Itoa(floatingIP))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d unassigned\n", floatingIP), out)

	out, err = runCommand(t, "floating-ip", "assign", strconv.Itoa(floatingIP), strconv.Itoa(serverID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d assigned to server %d\n", floatingIP, serverID), out)

	out, err = runCommand(t, "floating-ip", "describe", strconv.Itoa(floatingIP))
	require.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv4
Name:\s+test-floating-ip
Description:\s+-
Created:.*?
IP:\s+(?:[0-9]{1,3}\.){3}[0-9]{1,3}
Blocked:\s+no
Home Location:\s+[a-z]{3}[0-9]*
Server:
\s+ID:\s+[0-9]+
\s+Name:\s+test-server
DNS:
.*?
Protection:
\s+Delete:\s+no
Labels:
\s+No labels
`, out)

	out, err = runCommand(t, "floating-ip", "list", "-o", "columns=server", "-o", "noheader")
	require.NoError(t, err)
	assert.Equal(t, "test-server\n", out)

	out, err = runCommand(t, "floating-ip", "delete", strconv.Itoa(floatingIP))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIP), out)

	out, err = runCommand(t, "server", "delete", strconv.Itoa(serverID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Server %d deleted\n", serverID), out)
}
