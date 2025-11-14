//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestFloatingIP(t *testing.T) {
	t.Parallel()

	t.Run("ipv4", func(t *testing.T) {
		floatingIPName := withSuffix("test-floating-ip")
		_, err := createFloatingIP(t, floatingIPName, "")
		require.EqualError(t, err, "type is required")

		_, err = createFloatingIP(t, floatingIPName, "ipv4")
		require.EqualError(t, err, "one of --home-location or --server is required")

		_, err = createFloatingIP(t, floatingIPName, "ipv4", "--server", "non-existing-server")
		require.EqualError(t, err, "Server not found: non-existing-server")

		floatingIPId, err := createFloatingIP(t, floatingIPName, "ipv4", "--home-location", TestLocationName)
		require.NoError(t, err)

		t.Run("labels", func(t *testing.T) {
			t.Run("add-label-non-existing-floating-ip", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "add-label", "non-existing-floating-ip", "foo=bar")
				require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
				assert.Empty(t, out)
			})

			t.Run("add-label", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "add-label", strconv.FormatInt(floatingIPId, 10), "foo=bar")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Label(s) foo added to Floating IP %d\n", floatingIPId), out)
			})
		})

		t.Run("update-name", func(t *testing.T) {
			floatingIPName = withSuffix("new-test-floating-ip")
			out, err := runCommand(t, "floating-ip", "update", strconv.FormatInt(floatingIPId, 10), "--name", floatingIPName, "--description", "Some description")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Floating IP %d updated\n", floatingIPId), out)
		})

		t.Run("set-rnds", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "set-rdns", strconv.FormatInt(floatingIPId, 10), "--hostname", "s1.example.com")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)
		})

		t.Run("unassign-non-existing", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "unassign", "non-existing-floating-ip")
			require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
			assert.Empty(t, out)
		})

		t.Run("assign", func(t *testing.T) {
			t.Run("non-existing-ip", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "assign", "non-existing-floating-ip", "non-existing-server")
				require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
				assert.Empty(t, out)
			})

			t.Run("non-existing-server", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "assign", strconv.FormatInt(floatingIPId, 10), "non-existing-server")
				require.EqualError(t, err, "Server not found: non-existing-server")
				assert.Empty(t, out)
			})
		})

		t.Run("enable-protection", func(t *testing.T) {
			t.Run("unknown-protection-level", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "enable-protection", strconv.FormatInt(floatingIPId, 10), "unknown-protection-level")
				require.EqualError(t, err, "unknown protection level: unknown-protection-level")
				assert.Empty(t, out)
			})

			t.Run("non-existing-floating-ip", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "enable-protection", "non-existing-floating-ip", "delete")
				require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
				assert.Empty(t, out)
			})

			t.Run("enable-delete-protection", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "enable-protection", strconv.FormatInt(floatingIPId, 10), "delete")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Resource protection enabled for Floating IP %d\n", floatingIPId), out)
			})
		})

		var ipStr string

		t.Run("describe", func(t *testing.T) {
			t.Run("format", func(t *testing.T) {
				var err error
				ipStr, err = runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10), "--output", "format={{.IP}}")
				require.NoError(t, err)
				ipStr = strings.TrimSpace(ipStr)
				assert.Regexp(t, NewRegex().Start().IPv4().End(), ipStr)
			})

			t.Run("normal", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10))
				require.NoError(t, err)
				assert.Regexp(t,
					NewRegex().Start().
						Lit("ID:").Whitespace().Int().Newline().
						Lit("Type:").Whitespace().Lit("ipv4").Newline().
						Lit("Name:").Whitespace().Raw(`new-test-floating-ip-[0-9a-f]{8}`).Newline().
						Lit("Description:").Whitespace().Lit("Some description").Newline().
						Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
						Lit("IP:").Whitespace().IPv4().Newline().
						Lit("Blocked:").Whitespace().Lit("no").Newline().
						Lit("Home Location:").Whitespace().LocationName().Newline().
						Newline().
						Lit("Server:").Newline().
						Lit("  Not assigned").Newline().
						Newline().
						Lit("DNS:").Newline().
						Lit("  ").IPv4().Lit(":").Whitespace().Lit("s1.example.com").Newline().
						Newline().
						Lit("Protection:").Newline().
						Lit("  Delete:").Whitespace().Lit("yes").Newline().
						Newline().
						Lit("Labels:").Newline().
						Lit("  foo:").Whitespace().Lit("bar").Newline().
						End(),
					out,
				)
			})
		})

		t.Run("list", func(t *testing.T) {
			t.Run("table", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "list", "--output", "columns=id,name,type,ip,dns,server,home,blocked,protection,labels,created,age")
				require.NoError(t, err)
				assert.Regexp(t,
					NewRegex().Start().
						SeparatedByWhitespace("ID", "NAME", "TYPE", "IP", "DNS", "SERVER", "HOME", "BLOCKED", "PROTECTION", "LABELS", "CREATED", "AGE").Newline().
						Int().Whitespace().
						Raw(`new-test-floating-ip-[0-9a-f]{8}`).Whitespace().
						Lit("ipv4").Whitespace().
						IPv4().Whitespace().
						Lit("s1.example.com").Whitespace().
						Lit("-").Whitespace().
						LocationName().Whitespace().
						Lit("no").Whitespace().
						Lit("delete").Whitespace().
						Lit("foo=bar").Whitespace().
						Datetime().Whitespace().
						Age().OptionalWhitespace().Newline().
						End(),
					out,
				)
			})

			t.Run("json", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "list", "-o=json")
				require.NoError(t, err)
				assert.JSONEq(t, fmt.Sprintf(`
[
  {
    "id": %d,
    "description": "Some description",
    "created": "<ignore-diff>",
    "ip": "%s",
    "type": "ipv4",
    "server": null,
    "dns_ptr": [
      {
        "ip": "%s",
        "dns_ptr": "s1.example.com"
      }
    ],
    "home_location": {
      "id": "<ignore-diff>",
      "name": "%s",
      "description": "<ignore-diff>",
      "country": "<ignore-diff>",
      "city": "<ignore-diff>",
      "latitude": "<ignore-diff>",
      "longitude": "<ignore-diff>",
      "network_zone": "<ignore-diff>"
    },
    "blocked": false,
    "protection": {
      "delete": true
    },
    "labels": {
      "foo": "bar"
    },
    "name": "%s"
  }
]
`, floatingIPId, ipStr, ipStr, TestLocationName, floatingIPName), out)
			})
		})

		t.Run("delete-protected", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIPId, 10))
			require.Error(t, err)
			assert.Regexp(t, `^Floating IP deletion is protected \(protected, [0-9a-f]+\)$`, err.Error())
			assert.Empty(t, out)
		})

		t.Run("disable-protection", func(t *testing.T) {
			t.Run("non-existing-floating-ip", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "disable-protection", "non-existing-floating-ip", "delete")
				require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
				assert.Empty(t, out)
			})

			t.Run("unknown-protection-level", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "disable-protection", strconv.FormatInt(floatingIPId, 10), "unknown-protection-level")
				require.EqualError(t, err, "unknown protection level: unknown-protection-level")
				assert.Empty(t, out)
			})

			t.Run("disable-delete-protection", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "disable-protection", strconv.FormatInt(floatingIPId, 10), "delete")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Resource protection disabled for Floating IP %d\n", floatingIPId), out)
			})
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIPId, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIPId), out)
		})
	})

	t.Run("ipv6", func(t *testing.T) {
		floatingIPName := withSuffix("test-floating-ipv6")
		floatingIPId, err := createFloatingIP(t, floatingIPName, "ipv6", "--home-location", TestLocationName)
		require.NoError(t, err)

		var ipStr string

		t.Run("describe", func(t *testing.T) {
			t.Run("format", func(t *testing.T) {
				var err error
				ipStr, err = runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10), "--output", "format={{.IP}}")
				require.NoError(t, err)
				ipStr = strings.TrimSpace(ipStr)
				assert.NotNil(t, net.ParseIP(ipStr))
			})

			t.Run("normal", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10))
				require.NoError(t, err)
				assert.Regexp(t,
					NewRegex().Start().
						Lit("ID:").Whitespace().Int().Newline().
						Lit("Type:").Whitespace().Lit("ipv6").Newline().
						Lit("Name:").Whitespace().Raw(`test-floating-ipv6-[0-9a-f]{8}`).Newline().
						Lit("Description:").Whitespace().Lit("-").Newline().
						Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
						Lit("IP:").Whitespace().IPv6().Lit("/64").Newline().
						Lit("Blocked:").Whitespace().Lit("no").Newline().
						Lit("Home Location:").Whitespace().LocationName().Newline().
						Newline().
						Lit("Server:").Newline().
						Lit("  Not assigned").Newline().
						Newline().
						Lit("DNS:").Newline().
						Lit("  No reverse DNS entries").Newline().
						Newline().
						Lit("Protection:").Newline().
						Lit("  Delete:").Whitespace().Lit("no").Newline().
						Newline().
						Lit("Labels:").Newline().
						Lit("  No labels").Newline().
						End(),
					out,
				)
			})
		})

		t.Run("set-rdns", func(t *testing.T) {
			t.Run("1", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "set-rdns", strconv.FormatInt(floatingIPId, 10), "--ip", ipStr+"1", "--hostname", "s1.example.com")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)
			})

			t.Run("2", func(t *testing.T) {
				out, err := runCommand(t, "floating-ip", "set-rdns", strconv.FormatInt(floatingIPId, 10), "--ip", ipStr+"2", "--hostname", "s2.example.com")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)
			})
		})

		t.Run("list", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "list", "-o", "columns=ip,dns")
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("IP", "DNS").OptionalWhitespace().Newline().
					Lit(ipStr).Lit("/64").Whitespace().Lit("2 entries").OptionalWhitespace().Newline().
					End(),
				out,
			)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIPId, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIPId), out)
		})
	})
}

func createFloatingIP(t *testing.T, name, ipType string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.FloatingIP.Delete(context.Background(), &hcloud.FloatingIP{Name: name})
	})

	out, err := runCommand(t, append([]string{"floating-ip", "create", "--name", name, "--type", ipType}, args...)...)
	if err != nil {
		return 0, err
	}

	firstLine := strings.Split(out, "\n")[0]
	if !assert.Regexp(t, `^Floating IP [0-9]+ created$`, firstLine) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[12:len(firstLine)-8], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.FloatingIP.Delete(context.Background(), &hcloud.FloatingIP{ID: id})
	})
	return id, nil
}
