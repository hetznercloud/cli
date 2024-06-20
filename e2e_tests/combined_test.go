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
	assert.Equal(t, fmt.Sprintf("Firewall %d applied\n", firewallId), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallId))
	assert.Regexp(t, `^firewall with ID [0-9]+ is still in use \(resource_in_use, [0-9a-f]+\)$`, err.Error())
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "test-server", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d applied\n", firewallId), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallId), out)

	out, err = runCommand(t, "server", "delete", strconv.Itoa(serverId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Server %d deleted\n", serverId), out)
}
