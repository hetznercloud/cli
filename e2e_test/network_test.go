package e2e_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"
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

	networkID, err := strconv.Atoi(out[8 : len(out)-9])
	assert.NoError(t, err)

	out, err = runCommand(t, "network", "create", "--name", "test-network", "--ip-range", "10.0.1.0/24")
	assert.Empty(t, out)
	assert.Error(t, err)
	assert.Regexp(t, `^name is already used \(uniqueness_error, [0-9a-f]+\)$`, err.Error())

	out, err = runCommand(t, "network", "enable-protection", strconv.Itoa(networkID), "non-existing-protection")
	assert.EqualError(t, err, "unknown protection level: non-existing-protection")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "enable-protection", "non-existing-network", "delete")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "enable-protection", strconv.Itoa(networkID), "delete")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Resource protection enabled for network %d\n", networkID), out)

	out, err = runCommand(t, "network", "list", "-o=columns=servers,ip_range,labels,protection,created,age")
	assert.NoError(t, err)
	assert.Regexp(t, `SERVERS +IP RANGE +LABELS +PROTECTION +CREATED +AGE
0 servers +10\.0\.0\.0/24 +delete .*? (?:just now|[0-9]+s)
`, out)

	out, err = runCommand(t, "network", "change-ip-range", "--ip-range", "10.0.2.0/16", "non-existing-network")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "change-ip-range", "--ip-range", "10.0.2.0/16", "test-network")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("IP range of network %d changed\n", networkID), out)

	out, err = runCommand(t, "network", "add-label", "non-existing-network", "foo=bar")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "add-label", "test-network", "foo=bar")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to Network %d\n", networkID), out)

	out, err = runCommand(t, "network", "add-label", "test-network", "baz=qux")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to Network %d\n", networkID), out)

	out, err = runCommand(t, "network", "remove-label", "test-network", "baz")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Network %d\n", networkID), out)

	out, err = runCommand(t, "network", "update", "test-network", "--name", "new-test-network")
	assert.NoError(t, err)
	assert.Equal(t, "Network test-network updated\n", out)

	out, err = runCommand(t, "network", "delete", strconv.Itoa(networkID))
	assert.Empty(t, out)
	assert.Regexp(t, `^network is delete protected \(protected, [0-9a-f]+\)$`, err.Error())

	out, err = runCommand(t, "network", "add-subnet", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.16.0/24", "non-existing-network")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "add-subnet", "--type", "vswitch", "--vswitch-id", "42", "--network-zone", "eu-central", "--ip-range", "10.0.17.0/24", strconv.Itoa(networkID))
	assert.Empty(t, out)
	assert.Regexp(t, `^vswitch not found \(service_error, [0-9a-f]+\)$`, err.Error())

	out, err = runCommand(t, "network", "add-subnet", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.16.0/24", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Subnet added to network %d\n", networkID), out)

	out, err = runCommand(t, "network", "add-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", "non-existing-network")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "add-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Route added to network %d\n", networkID), out)

	out, err = runCommand(t, "network", "expose-routes-to-vswitch", "non-existing-network")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "expose-routes-to-vswitch", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, "Exposing routes to connected vSwitch of network new-test-network enabled\n", out)

	out, err = runCommand(t, "network", "describe", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Regexp(t, `^ID:\s+[0-9]+
Name:\s+new-test-network
Created:\s+.*?
IP Range:\s+10\.0\.0\.0\/16
Expose Routes to vSwitch: yes
Subnets:
\s+- Type:\s+cloud
\s+Network Zone:\s+eu-central
\s+IP Range:\s+10\.0\.16\.0\/24
\s+Gateway:\s+10\.0\.0\.1
Routes:
\s+- Destination:\s+10\.100\.1\.0\/24
\s+Gateway:\s+10\.0\.1\.1
Protection:
\s+Delete:\s+yes
Labels:
\s+foo: bar
$`, out)

	out, err = runCommand(t, "network", "list", "-o=json")
	assert.NoError(t, err)
	assertjson.Equal(t, []byte(fmt.Sprintf(`
[
  {
    "id": %d,
    "name": "new-test-network",
    "created": "<ignore-diff>",
    "ip_range": "10.0.0.0/16",
    "subnets": [
      {
        "type": "cloud",
        "ip_range": "10.0.16.0/24",
        "network_zone": "eu-central",
        "gateway": "10.0.0.1"
      }
    ],
    "routes": [
      {
        "destination": "10.100.1.0/24",
        "gateway": "10.0.1.1"
      }
    ],
    "servers": [],
    "protection": {
      "delete": true
    },
    "labels": {
      "foo": "bar"
    },
    "expose_routes_to_vswitch": true
  }
]
`, networkID)), []byte(out))

	out, err = runCommand(t, "network", "remove-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", "non-existing-network")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "remove-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Route removed from network %d\n", networkID), out)

	out, err = runCommand(t, "network", "remove-subnet", "--ip-range", "10.0.16.0/24", "non-existing-network")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "remove-subnet", "--ip-range", "10.0.16.0/24", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Subnet 10.0.16.0/24 removed from network %d\n", networkID), out)

	out, err = runCommand(t, "network", "disable-protection", "non-existing-network", "delete")
	assert.EqualError(t, err, "network not found: non-existing-network")
	assert.Empty(t, out)

	out, err = runCommand(t, "network", "disable-protection", strconv.Itoa(networkID), "delete")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Resource protection disabled for network %d\n", networkID), out)

	out, err = runCommand(t, "network", "remove-label", strconv.Itoa(networkID), "foo")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo removed from Network %d\n", networkID), out)

	out, err = runCommand(t, "network", "expose-routes-to-vswitch", "--disable", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, "Exposing routes to connected vSwitch of network new-test-network disabled\n", out)

	out, err = runCommand(t, "network", "describe", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Regexp(t, `^ID:\s+[0-9]+
Name:\s+new-test-network
Created:\s+.*?
IP Range:\s+10\.0\.0\.0\/16
Expose Routes to vSwitch: no
Subnets:
\s+No subnets
Routes:
\s+No routes
Protection:
\s+Delete:\s+no
Labels:
\s+No labels
$`, out)

	out, err = runCommand(t, "network", "delete", strconv.Itoa(networkID))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Network %d deleted\n", networkID), out)
}
