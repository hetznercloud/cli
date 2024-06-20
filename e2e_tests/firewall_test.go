package e2e_tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirewall(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "firewall", "create", "--name", "test-firewall", "--rules-file", "rules_file.json")
	assert.NoError(t, err)
	if !assert.Regexp(t, `^Firewall [0-9]+ created\n$`, out) {
		// firewall was not created (properly), so there's no need to test it
		return
	}

	firewallId, err := strconv.Atoi(out[9 : len(out)-9])
	assert.NoError(t, err)

	out, err = runCommand(t, "firewall", "add-label", "non-existing-firewall", "foo=bar")
	assert.EqualError(t, err, "firewall not found: non-existing-firewall")
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "add-label", strconv.Itoa(firewallId), "foo=bar")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to firewall %d\n", firewallId), out)

	out, err = runCommand(t, "firewall", "add-label", strconv.Itoa(firewallId), "baz=qux")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to firewall %d\n", firewallId), out)

	out, err = runCommand(t, "firewall", "update", strconv.Itoa(firewallId), "--name", "new-test-firewall")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d updated\n", firewallId), out)

	out, err = runCommand(t, "firewall", "remove-label", "new-test-firewall", "baz")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from firewall %d\n", firewallId), out)

	out, err = runCommand(t, "firewall", "add-rule", "non-existing-firewall", "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100")
	assert.EqualError(t, err, "Firewall not found: non-existing-firewall")
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallId), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--description", "Some random description")
	assert.EqualError(t, err, "port is required (--port)")
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallId), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100", "--description", "Some random description")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallId), out)

	out, err = runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallId), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp", "--port", "12345")
	assert.EqualError(t, err, "port is not allowed for this protocol")
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallId), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallId), out)

	out, err = runCommand(t, "firewall", "describe", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Name:\s+new-test-firewall
Created:\s+.*?
Labels:
\s+foo: bar
Rules:
\s+- Direction:\s+in
\s+Description:\s+Allow port 80
\s+Protocol:\s+tcp
\s+Port:\s+80
\s+Source IPs:
\s+28\.239\.13\.1\/32
\s+28\.239\.14\.0\/24
\s+ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b\/128
\s+- Direction:\s+in
\s+Description:\s+Allow port 443
\s+Protocol:\s+tcp
\s+Port:\s+443
\s+Source IPs:
\s+0\.0\.0\.0\/0
\s+::\/0
\s+- Direction:\s+out
\s+Protocol:\s+tcp
\s+Port:\s+80
\s+Destination IPs:
\s+28\.239\.13\.1\/32
\s+28\.239\.14\.0\/24
\s+ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b\/128
\s+- Direction:\s+in
\s+Description:\s+Some random description
\s+Protocol:\s+tcp
\s+Port:\s+9100
\s+Source IPs:
\s+10\.0\.0\.0\/24
\s+- Direction:\s+out
\s+Protocol:\s+icmp
\s+Destination IPs:
\s+192\.168\.1\.0\/24
Applied To:
\s+Not applied
`, out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallId))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallId), out)
}
