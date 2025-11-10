//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestFirewall(t *testing.T) {
	t.Parallel()

	firewallName := withSuffix("test-firewall")
	firewallID, err := createFirewall(t, firewallName, "--rules-file", "rules_file.json")
	require.NoError(t, err)

	t.Run("add-label", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-label", "non-existing-firewall", "foo=bar")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("1", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-label", strconv.FormatInt(firewallID, 10), "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to Firewall %d\n", firewallID), out)
		})

		t.Run("2", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-label", strconv.FormatInt(firewallID, 10), "baz=qux")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz added to Firewall %d\n", firewallID), out)
		})
	})

	t.Run("update-name", func(t *testing.T) {
		firewallName = withSuffix("new-test-firewall")

		out, err := runCommand(t, "firewall", "update", strconv.FormatInt(firewallID, 10), "--name", firewallName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Firewall %d updated\n", firewallID), out)
	})

	t.Run("remove-label", func(t *testing.T) {
		out, err := runCommand(t, "firewall", "remove-label", firewallName, "baz")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Firewall %d\n", firewallID), out)
	})

	t.Run("add-rule", func(t *testing.T) {
		t.Run("missing-args", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, `required flag(s) "direction", "protocol" not set`)
			assert.Empty(t, out)
		})

		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", "non-existing-firewall", "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("missing-port", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--description", "Some random description")
			require.EqualError(t, err, "port is required (--port)")
			assert.Empty(t, out)
		})

		t.Run("port-not-allowed", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp", "--port", "12345")
			require.EqualError(t, err, "port is not allowed for this protocol")
			assert.Empty(t, out)
		})

		t.Run("invalid-direction", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "foo", "--destination-ips", "192.168.1.0/24", "--protocol", "tcp", "--port", "12345")
			require.EqualError(t, err, "invalid direction: foo")
			assert.Empty(t, out)
		})

		t.Run("invalid-protocol", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "abc", "--port", "12345")
			require.EqualError(t, err, "invalid protocol: abc")
			assert.Empty(t, out)
		})

		t.Run("tcp-in", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100", "--description", "Some random description")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("icmp-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("invalid-ip-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "out", "--destination-ips", "invalid-ip", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "destination error on index 0: invalid CIDR address: invalid-ip")
			assert.Empty(t, out)
		})

		t.Run("invalid-ip-range-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "add-rule", strconv.FormatInt(firewallID, 10), "--direction", "in", "--source-ips", "10.1.2.3/8", "--protocol", "tcp", "--port", "9100")
			require.EqualError(t, err, "source ips error on index 0: 10.1.2.3/8 is not the start of the cidr block 10.0.0.0/8")
			assert.Empty(t, out)
		})
	})

	t.Run("apply-to-resource", func(t *testing.T) {
		t.Run("unknown-type", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "non-existing-type", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "unknown type non-existing-type")
			assert.Empty(t, out)
		})

		t.Run("missing-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "type server need a --server specific")
			assert.Empty(t, out)
		})

		t.Run("missing-label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "label_selector", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "type label_selector need a --label-selector specific")
			assert.Empty(t, out)
		})

		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", "non-existing-server", "non-existing-firewall")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("unknown-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", "non-existing-server", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "Server not found: non-existing-server")
			assert.Empty(t, out)
		})

		t.Run("label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "label_selector", "--label-selector", "foo=bar", strconv.FormatInt(firewallID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall %d applied to resource\n", firewallID), out)
		})
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "firewall", "describe", strconv.FormatInt(firewallID, 10))
		require.NoError(t, err)
		assert.Regexp(t,
			NewRegex().Start().
				Lit("ID:").Whitespace().Int().Newline().
				Lit("Name:").Whitespace().Raw(`new-test-firewall-[0-9a-f]{8}`).Newline().
				Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
				Newline().
				Lit("Labels:").Newline().
				Lit("  foo:").Whitespace().Lit("bar").Newline().
				Newline().
				Lit("Rules:").Newline().
				Lit("  - Direction:").Whitespace().Lit("in").Newline().
				Lit("    Description:").Whitespace().Lit("Allow port 80").Newline().
				Lit("    Protocol:").Whitespace().Lit("tcp").Newline().
				Lit("    Port:").Whitespace().Lit("80").Newline().
				Lit("    Source IPs:").Newline().
				Whitespace().Lit("28.239.13.1/32").Newline().
				Whitespace().Lit("28.239.14.0/24").Newline().
				Whitespace().Lit("ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b/128").Newline().
				Newline().
				Lit("  - Direction:").Whitespace().Lit("in").Newline().
				Lit("    Description:").Whitespace().Lit("Allow port 443").Newline().
				Lit("    Protocol:").Whitespace().Lit("tcp").Newline().
				Lit("    Port:").Whitespace().Lit("443").Newline().
				Lit("    Source IPs:").Newline().
				Whitespace().Lit("0.0.0.0/0").Newline().
				Whitespace().Lit("::/0").Newline().
				Newline().
				Lit("  - Direction:").Whitespace().Lit("in").Newline().
				Lit("    Description:").Whitespace().Lit("Some random description").Newline().
				Lit("    Protocol:").Whitespace().Lit("tcp").Newline().
				Lit("    Port:").Whitespace().Lit("9100").Newline().
				Lit("    Source IPs:").Newline().
				Whitespace().Lit("10.0.0.0/24").Newline().
				Newline().
				Lit("  - Direction:").Whitespace().Lit("out").Newline().
				Lit("    Protocol:").Whitespace().Lit("tcp").Newline().
				Lit("    Port:").Whitespace().Lit("80").Newline().
				Lit("    Destination IPs:").Newline().
				Whitespace().Lit("28.239.13.1/32").Newline().
				Whitespace().Lit("28.239.14.0/24").Newline().
				Whitespace().Lit("ff21:1eac:9a3b:ee58:5ca:990c:8bc9:c03b/128").Newline().
				Newline().
				Lit("  - Direction:").Whitespace().Lit("out").Newline().
				Lit("    Protocol:").Whitespace().Lit("icmp").Newline().
				Lit("    Destination IPs:").Newline().
				Whitespace().Lit("192.168.1.0/24").Newline().
				Newline().
				Lit("Applied To:").Newline().
				Lit("  - Type:").Whitespace().Lit("label_selector").Newline().
				Lit("    Label Selector:").Whitespace().Lit("foo=bar").Newline().
				End(),
			out,
		)
	})

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "list", "--output", "columns=id,name,rules_count,applied_to_count,labels,created,age")
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "NAME", "RULES COUNT", "APPLIED TO COUNT", "LABELS", "CREATED", "AGE").Newline().
					Int().Whitespace().
					Raw(`new-test-firewall-[0-9a-f]{8}`).Whitespace().
					Lit("5 Rules").Whitespace().
					Lit("0 Servers | 1 Label Selector").Whitespace().
					Lit("foo=bar").Whitespace().
					Datetime().Whitespace().
					Age().OptionalWhitespace().Newline().
					End(),
				out,
			)
		})
	})

	t.Run("remove-from-resource", func(t *testing.T) {
		t.Run("unknown-type", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "non-existing-type", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "unknown type non-existing-type")
			assert.Empty(t, out)
		})

		t.Run("missing-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "type server need a --server specific")
			assert.Empty(t, out)
		})

		t.Run("missing-label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "label_selector", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "type label_selector need a --label-selector specific")
			assert.Empty(t, out)
		})

		t.Run("unknown-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "non-existing-server", "non-existing-firewall")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("unknown-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", "non-existing-server", strconv.FormatInt(firewallID, 10))
			require.EqualError(t, err, "Server not found: non-existing-server")
			assert.Empty(t, out)
		})

		t.Run("label-selector", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "label_selector", "--label-selector", "foo=bar", strconv.FormatInt(firewallID, 10))
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
			out, err := runCommand(t, "firewall", "delete-rule", strconv.FormatInt(firewallID, 10), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp")
			require.EqualError(t, err, "port is required (--port)")
			assert.Empty(t, out)
		})

		t.Run("port-not-allowed", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.FormatInt(firewallID, 10), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp", "--port", "12345")
			require.EqualError(t, err, "port is not allowed for this protocol")
			assert.Empty(t, out)
		})

		t.Run("tcp-in", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.FormatInt(firewallID, 10), "--direction", "in", "--source-ips", "10.0.0.0/24", "--protocol", "tcp", "--port", "9100", "--description", "Some random description")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("icmp-out", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.FormatInt(firewallID, 10), "--direction", "out", "--destination-ips", "192.168.1.0/24", "--protocol", "icmp")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})

		t.Run("non-existing-rule", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete-rule", strconv.FormatInt(firewallID, 10), "--direction", "in", "--source-ips", "123.123.123.123/32", "--port", "1234", "--protocol", "tcp")
			require.EqualError(t, err, fmt.Sprintf("the specified rule was not found in the ruleset of Firewall %d", firewallID))
			assert.Empty(t, out)
		})
	})

	t.Run("replace-rules", func(t *testing.T) {
		t.Run("non-existing-firewall", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "replace-rules", "non-existing-firewall", "--rules-file", "rules_file.json")
			require.EqualError(t, err, "Firewall not found: non-existing-firewall")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "replace-rules", strconv.FormatInt(firewallID, 10), "--rules-file", "rules_file.json")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall Rules for Firewall %d updated\n", firewallID), out)
		})
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "firewall", "delete", strconv.FormatInt(firewallID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Firewall %d deleted\n", firewallID), out)
	})
}

func createFirewall(t *testing.T, name string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.Firewall.Delete(context.Background(), &hcloud.Firewall{Name: name})
	})

	out, err := runCommand(t, append([]string{"firewall", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^Firewall [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[9:len(out)-9], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.Firewall.Delete(context.Background(), &hcloud.Firewall{ID: id})
	})
	return id, nil
}
