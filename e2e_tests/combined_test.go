package e2e_tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombined(t *testing.T) {
	// combined tests combine multiple resources and can thus not be run in parallel
	serverId := createServer(t, "test-server", "cpx11", "ubuntu-24.04")

	firewallId, err := createFirewall(t, "test-firewall")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", "test-server", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d applied to resource\n", firewallId), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallId))
	assert.Regexp(t, `^firewall with ID [0-9]+ is still in use \(resource_in_use, [0-9a-f]+\)$`, err.Error())
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "test-server", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d removed from resource\n", firewallId), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallId), out)

	floatingIP, err := createFloatingIP(t, "test-floating-ip", "ipv4", "--server", strconv.Itoa(serverId))
	if err != nil {
		t.Fatal(err)
	}

	out, err = runCommand(t, "floating-ip", "unassign", strconv.Itoa(floatingIP))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d unassigned\n", floatingIP), out)

	out, err = runCommand(t, "floating-ip", "assign", strconv.Itoa(floatingIP), strconv.Itoa(serverId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d assigned to server %d\n", floatingIP, serverId), out)

	out, err = runCommand(t, "floating-ip", "describe", strconv.Itoa(floatingIP))
	assert.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv4
Name:\s+test-floating-ip
Description:\s+-
Created:.*?
IP:\s+(?:[0-9]{1,3}\.){3}[0-9]{1,3}
Blocked:\s+no
Home Location:\s+[a-z]{3}[0-9]+
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
	assert.NoError(t, err)
	assert.Equal(t, "test-server\n", out)

	out, err = runCommand(t, "floating-ip", "delete", strconv.Itoa(floatingIP))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIP), out)

	out, err = runCommand(t, "server", "delete", strconv.Itoa(serverId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Server %d deleted\n", serverId), out)
}
