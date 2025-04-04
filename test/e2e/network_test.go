//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/assertjson"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestNetwork(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "network", "create")
	assert.Empty(t, out)
	require.EqualError(t, err, `required flag(s) "ip-range", "name" not set`)

	networkName := withSuffix("test-network")
	networkID, err := createNetwork(t, networkName, "--ip-range", "10.0.0.0/24")
	require.NoError(t, err)

	_, err = createNetwork(t, networkName, "--ip-range", "10.0.1.0/24")
	require.Error(t, err)
	assert.Regexp(t, `^name is already used \(uniqueness_error, [0-9a-f]+\)$`, err.Error())

	t.Run("enable-protection", func(t *testing.T) {
		t.Run("non-existing-protection", func(t *testing.T) {
			out, err := runCommand(t, "network", "enable-protection", strconv.FormatInt(networkID, 10), "non-existing-protection")
			require.EqualError(t, err, "unknown protection level: non-existing-protection")
			assert.Empty(t, out)
		})

		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "enable-protection", "non-existing-network", "delete")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "network", "enable-protection", strconv.FormatInt(networkID, 10), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection enabled for Network %d\n", networkID), out)
		})
	})

	t.Run("list", func(t *testing.T) {
		out, err := runCommand(t, "network", "list", "-o=columns=servers,ip_range,labels,protection,created,age")
		require.NoError(t, err)
		assert.Regexp(t,
			NewRegex().Start().
				SeparatedByWhitespace("SERVERS", "IP", "RANGE", "LABELS", "PROTECTION", "CREATED", "AGE").OptionalWhitespace().Newline().
				Lit("0 servers").Whitespace().
				Lit("10.0.0.0/24").Whitespace().
				Lit("delete").Whitespace().
				UnixDate().Whitespace().
				Age().OptionalWhitespace().Newline().End(),
			out,
		)
	})

	t.Run("change-ip-range", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "change-ip-range", "--ip-range", "10.0.2.0/16", "non-existing-network")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "change-ip-range", "--ip-range", "10.0.2.0/16", networkName)
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("IP range of Network %d changed\n", networkID), out)
		})
	})

	t.Run("labels", func(t *testing.T) {
		t.Run("add", func(t *testing.T) {
			t.Run("non-existing-network", func(t *testing.T) {
				out, err := runCommand(t, "network", "add-label", "non-existing-network", "foo=bar")
				require.EqualError(t, err, "Network not found: non-existing-network")
				assert.Empty(t, out)
			})

			t.Run("1", func(t *testing.T) {
				out, err := runCommand(t, "network", "add-label", networkName, "foo=bar")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Label(s) foo added to Network %d\n", networkID), out)
			})

			t.Run("2", func(t *testing.T) {
				out, err := runCommand(t, "network", "add-label", networkName, "baz=qux")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Label(s) baz added to Network %d\n", networkID), out)
			})
		})

		t.Run("remove", func(t *testing.T) {
			out, err := runCommand(t, "network", "remove-label", networkName, "baz")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Network %d\n", networkID), out)
		})
	})

	oldNetworkName := networkName
	networkName = withSuffix("new-test-network")

	t.Run("update-name", func(t *testing.T) {
		out, err := runCommand(t, "network", "update", oldNetworkName, "--name", networkName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Network %s updated\n", oldNetworkName), out)
	})

	t.Run("delete-protected", func(t *testing.T) {
		out, err := runCommand(t, "network", "delete", strconv.FormatInt(networkID, 10))
		assert.Empty(t, out)
		assert.Regexp(t, `^network is delete protected \(protected, [0-9a-f]+\)$`, err.Error())
	})

	t.Run("add-subnet", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "add-subnet", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.16.0/24", "non-existing-network")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("non-existing-vswitch", func(t *testing.T) {
			out, err := runCommand(t, "network", "add-subnet", "--type", "vswitch", "--vswitch-id", "42", "--network-zone", "eu-central", "--ip-range", "10.0.17.0/24", strconv.FormatInt(networkID, 10))
			assert.Empty(t, out)
			assert.Regexp(t, `^vswitch not found \(service_error, [0-9a-f]+\)$`, err.Error())
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "add-subnet", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.16.0/24", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Subnet added to Network %d\n", networkID), out)
		})
	})

	t.Run("add-route", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "add-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", "non-existing-network")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)

		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "add-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Route added to Network %d\n", networkID), out)
		})
	})

	t.Run("expose-routes-to-vswitch", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "expose-routes-to-vswitch", "non-existing-network")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "expose-routes-to-vswitch", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Exposing routes to connected vSwitch of Network %s enabled\n", networkName), out)
		})
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "network", "describe", strconv.FormatInt(networkID, 10))
		require.NoError(t, err)
		assert.Regexp(t, NewRegex().Start().
			Lit("ID:").Whitespace().Int().Newline().
			Lit("Name:").Whitespace().Raw(`new-test-network-[0-9a-f]{8}`).Newline().
			Lit("Created:").Whitespace().UnixDate().Lit(" (").HumanizeTime().Lit(")").Newline().
			Lit("IP Range:").Whitespace().Lit("10.0.0.0/16").Newline().
			Lit("Expose Routes to vSwitch:").Whitespace().Lit("yes").Newline().
			Lit("Subnets:").Newline().
			Lit("  - Type:").Whitespace().Lit("cloud").Newline().
			Lit("    Network Zone:").Whitespace().OneOfLit("eu-central", "us-east", "us-west", "ap-southeast").Newline().
			Lit("    IP Range:").Whitespace().Lit("10.0.16.0/24").Newline().
			Lit("    Gateway:").Whitespace().Lit("10.0.0.1").Newline().
			Lit("Routes:").Newline().
			Lit("  - Destination:").Whitespace().Lit("10.100.1.0/24").Newline().
			Lit("    Gateway:").Whitespace().Lit("10.0.1.1").Newline().
			Lit("Protection:").Newline().
			Lit("  Delete:").Whitespace().Lit("yes").Newline().
			Lit("Labels:").Newline().
			Lit("  foo: bar").Newline().
			End(),
			out,
		)
	})

	t.Run("list", func(t *testing.T) {
		out, err := runCommand(t, "network", "list", "-o=json")
		require.NoError(t, err)
		assertjson.Equal(t, []byte(fmt.Sprintf(`
[
  {
    "id": %d,
    "name": "%s",
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
    "load_balancers": [],
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
`, networkID, networkName)), []byte(out))
	})

	t.Run("remove-route", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "remove-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", "non-existing-network")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "remove-route", "--destination", "10.100.1.0/24", "--gateway", "10.0.1.1", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Route removed from Network %d\n", networkID), out)
		})
	})

	t.Run("remove-subnet", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "remove-subnet", "--ip-range", "10.0.16.0/24", "non-existing-network")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "remove-subnet", "--ip-range", "10.0.16.0/24", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Subnet 10.0.16.0/24 removed from Network %d\n", networkID), out)
		})
	})

	t.Run("disable-protection", func(t *testing.T) {
		t.Run("non-existing-network", func(t *testing.T) {
			out, err := runCommand(t, "network", "disable-protection", "non-existing-network", "delete")
			require.EqualError(t, err, "Network not found: non-existing-network")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "network", "disable-protection", strconv.FormatInt(networkID, 10), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection disabled for Network %d\n", networkID), out)
		})
	})

	t.Run("remove-label", func(t *testing.T) {
		out, err := runCommand(t, "network", "remove-label", strconv.FormatInt(networkID, 10), "foo")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Label(s) foo removed from Network %d\n", networkID), out)
	})

	t.Run("disable-expose-routes-to-vswitch", func(t *testing.T) {
		out, err := runCommand(t, "network", "expose-routes-to-vswitch", "--disable", strconv.FormatInt(networkID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Exposing routes to connected vSwitch of Network %s disabled\n", networkName), out)
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "network", "describe", strconv.FormatInt(networkID, 10))
		require.NoError(t, err)
		assert.Regexp(t, NewRegex().Start().
			Lit("ID:").Whitespace().Int().Newline().
			Lit("Name:").Whitespace().Raw(`new-test-network-[0-9a-f]{8}`).Newline().
			Lit("Created:").Whitespace().UnixDate().Lit(" (").HumanizeTime().Lit(")").Newline().
			Lit("IP Range:").Whitespace().Lit("10.0.0.0/16").Newline().
			Lit("Expose Routes to vSwitch:").Whitespace().Lit("no").Newline().
			Lit("Subnets:").Newline().
			Lit("  No subnets").Newline().
			Lit("Routes:").Newline().
			Lit("  No routes").Newline().
			Lit("Protection:").Newline().
			Lit("  Delete:").Whitespace().Lit("no").Newline().
			Lit("Labels:").Newline().
			Lit("  No labels").Newline().
			End(),
			out,
		)
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "network", "delete", strconv.FormatInt(networkID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Network %d deleted\n", networkID), out)
	})
}

func createNetwork(t *testing.T, name string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.Network.Delete(context.Background(), &hcloud.Network{Name: name})
	})

	out, err := runCommand(t, append([]string{"network", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^Network [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[8:len(out)-9], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.Network.Delete(context.Background(), &hcloud.Network{ID: id})
	})
	return id, nil
}
