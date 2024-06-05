package e2e_tests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "network", "create")
	assert.Empty(t, out)
	assert.EqualError(t, err, `required flag(s) "ip-range", "name" not set`)

	out, err = runCommand(t, "network", "create", "--name", "test-network", "--ip-range", "10.0.0.0/24")
	assert.NoError(t, err)
	if !assert.Regexp(t, `^Network [0-9]+ created\n$`, out) {
		// network was not created (properly), so there's no need to test it
		return
	}

	networkId, err := strconv.Atoi(out[8 : len(out)-9])
	assert.NoError(t, err)

	out, err = runCommand(t, "network", "create", "--name", "test-network", "--ip-range", "10.0.1.0/24")
	assert.Empty(t, out)
	assert.Error(t, err)
	assert.Regexp(t, `^name is already used \(uniqueness_error, [0-9a-f]{16}\)$`, err.Error())

	out, err = runCommand(t, "network", "list")
	assert.NoError(t, err)
	lines := strings.Split(out, "\n")
	if assert.Len(t, lines, 3) {
		assert.Regexp(t, "ID +NAME +IP +RANGE +SERVERS +AGE", lines[0])
		assert.Regexp(t, fmt.Sprintf(`^%d +test-network +10\.0\.0\.0/24 +0 server +(?:just now|[0-9]+s)$`, networkId), lines[1])
		assert.Empty(t, lines[2])
	}

	out, err = runCommand(t, "network", "change-ip-range", "test-network", "--ip-range", "10.0.2.0/16")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("IP range of network %d changed\n", networkId), out)

	out, err = runCommand(t, "network", "add-label", "test-network", "foo=bar")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to Network %d\n", networkId), out)

	out, err = runCommand(t, "network", "add-label", "test-network", "baz=qux")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to Network %d\n", networkId), out)

	out, err = runCommand(t, "network", "remove-label", "test-network", "baz")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Network %d\n", networkId), out)

	out, err = runCommand(t, "network", "update", "test-network", "--name", "new-test-network")
	assert.NoError(t, err)
	assert.Equal(t, "Network test-network updated\n", out)

	out, err = runCommand(t, "network", "enable-protection", strconv.Itoa(networkId), "delete")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Resource protection enabled for network %d\n", networkId), out)

	out, err = runCommand(t, "network", "delete", strconv.Itoa(networkId))
	assert.Empty(t, out)
	assert.Regexp(t, `^network deletion is protected \(protected, [0-9a-f]{16}\)$`, err.Error())

	out, err = runCommand(t, "network", "add-subnet", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.16.0/24", strconv.Itoa(networkId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Subnet added to network %d\n", networkId), out)

	out, err = runCommand(t, "network", "add-route", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.16.0/24", strconv.Itoa(networkId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Subnet added to network %d\n", networkId), out)

	out, err = runCommand(t, "network", "disable-protection", strconv.Itoa(networkId), "delete")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Resource protection disabled for network %d\n", networkId), out)

	out, err = runCommand(t, "network", "delete", strconv.Itoa(networkId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Network %d deleted\n", networkId), out)
}
