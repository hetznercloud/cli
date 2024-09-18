//go:build e2e

package e2e_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirewall(t *testing.T) {
	t.Parallel()

	firewallID, err := createFirewall(t, "test-firewall", "--rules-file", "rules_file.json")
	if err != nil {
		t.Fatal(err)
	}

	out, err := runCommand(t, "firewall", "add-label", "non-existing-firewall", "foo=bar")
	require.EqualError(t, err, "firewall not found: non-existing-firewall")
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "add-label", strconv.Itoa(firewallID), "foo=bar")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to firewall %d\n", firewallID), out)

	out, err = runCommand(t, "firewall", "add-label", strconv.Itoa(firewallID), "baz=qux")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to firewall %d\n", firewallID), out)

	out, err = runCommand(t, "firewall", "update", strconv.Itoa(firewallID), "--name", "new-test-firewall")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall %d updated\n", firewallID), out)

	out, err = runCommand(t, "firewall", "remove-label", "new-test-firewall", "baz")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from firewall %d\n", firewallID), out)

	t.Run("add-rule", func(t *testing.T) {
		t.Run("missing-args", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID))
			require.EqualError(t, err, `required flag(s) "direction", "protocol" not set`)
			assert.Empty(t, out)
		})

		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", "non-existing-firewall", "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("missing-port", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--description", "Some random description")
			require.EqualError(t, err, "port is required (--port)")
			assert.Empty(t, out)
		})

		t.Run("port-not-allowed", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp", "--port", "12345")
			require.EqualError(t, err, "port is not allowed for this protocol")
			assert.Empty(t, out)
		})

		t.Run("invalid-direction", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "foo", "--destination-ips", "192.168.1.0/24", "--protocol", "tcp", "--port", "12345")
			require.EqualError(t, err, "invalid direction: foo")
			assert.Empty(t, out)
		})

		t.Run("invalid-protocol", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "abc", "--port", "12345")
			require.EqualError(t, err, "invalid protocol: abc")
			assert.Empty(t, out)
		})

		t.Run("tcp-in", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100", "--description", "Some random description")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("icmp-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("invalid-ip-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "out", "--destination-ips", "invalid-ip", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "destination error on index 0: invalid CIDR address: invalid-ip")
			assert.Empty(t, out)
		})

		t.Run("invalid-ip-range-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.Itoa(firewallID), "--direction", "in", "--source-ips", "10.1.2.3/8", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "source ips error on index 0: 10.1.2.3/8 is not the start of the cidr block 10.0.0.0/8")
			assert.Empty(t, out)
		})
	})

	t.Run("apply-to-resource", func(t *testing.T) {
		t.Run("unknown-type", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "non-existing-type", strconv.Itoa(firewallID))
			require.EqualError(t, err, "unknown type non-existing-type")
			assert.Empty(t, out)
		})
		t.Run("missing-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", strconv.Itoa(firewallID))
			require.EqualError(t, err, "type server need a --server specific")
			assert.Empty(t, out)
		})
		t.Run("missing-label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "label_selector", strconv.Itoa(firewallID))
			require.EqualError(t, err, "type label_selector need a --label-selector specific")
			assert.Empty(t, out)
		})
		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", "non-existing-server", "non-existing-firewall")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})
		t.Run("unknown-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", "non-existing-server", strconv.Itoa(firewallID))
			require.EqualError(t, err, "Server not found: non-existing-server")
			assert.Empty(t, out)
		})
		t.Run("label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "label_selector", "--label-selector", "foo=bar", strconv.Itoa(firewallID))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall %d applied to resource\n", firewallID), out)
		})
	})

	out, err = runCommand(t, "firewall", "describe", strconv.Itoa(firewallID))
	require.NoError(t, err)
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
\s+- Type:\s+label_selector
\s+Label Selector:\s+foo=bar
`, out)

	t.Run("remove-from-resource", func(t *testing.T) {
		t.Run("unknown-type", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "non-existing-type", strconv.Itoa(firewallID))
			require.EqualError(t, err, "unknown type non-existing-type")
			assert.Empty(t, out)
		})
		t.Run("missing-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", strconv.Itoa(firewallID))
			require.EqualError(t, err, "type server need a --server specific")
			assert.Empty(t, out)
		})
		t.Run("missing-label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "label_selector", strconv.Itoa(firewallID))
			require.EqualError(t, err, "type label_selector need a --label-selector specific")
			assert.Empty(t, out)
		})
		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "non-existing-server", "non-existing-firewall")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})
		t.Run("unknown-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "non-existing-server", strconv.Itoa(firewallID))
			require.EqualError(t, err, "Server not found: non-existing-server")
			assert.Empty(t, out)
		})
		t.Run("label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "label_selector", "--label-selector", "foo=bar", strconv.Itoa(firewallID))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall %d removed from resource\n", firewallID), out)
		})
	})

	t.Run("delete-rule", func(t *testing.T) {
		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", "non-existing-firewall", "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("missing-port", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.Itoa(firewallID), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp")
			require.EqualError(t, err, "port is required (--port)")
			assert.Empty(t, out)
		})

		t.Run("port-not-allowed", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.Itoa(firewallID), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp", "--port", "12345")
			require.EqualError(t, err, "port is not allowed for this protocol")
			assert.Empty(t, out)
		})

		t.Run("tcp-in", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.Itoa(firewallID), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100", "--description", "Some random description")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("icmp-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.Itoa(firewallID), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("non-existing-rule", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.Itoa(firewallID), "--direction", "in", "--source-ips", "123.123.123.123/32", "--port", "1234", "--protocol", "tcp")
			require.EqualError(t, err, fmt.Sprintf("the specified rule was not found in the ruleset of Firewall %d", firewallID))
			assert.Empty(t, out)
		})
	})

	out, err = runCommand(t, "firewall", "replace-rules", "non-existing-firewall", "--rules-file", "rules_file.json")
	require.EqualError(t, err, "Firewall not found: non-existing-firewall")
	assert.Empty(t, out)

	out, err = runCommand(t, "firewall", "replace-rules", strconv.Itoa(firewallID), "--rules-file", "rules_file.json")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)

	out, err = runCommand(t, "firewall", "delete", strconv.Itoa(firewallID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallID), out)
}

func createFirewall(t *testing.T, name string, args ...string) (int, error) {
	out, err := runCommand(t, append([]string{"firewall", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^Firewall [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	firewallID, err := strconv.Atoi(out[9 : len(out)-9])
	if err != nil {
		return 0, err
	}
	return firewallID, nil
}
