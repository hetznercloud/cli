//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLoadBalancer(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "load-balancer", "create")
	assert.Empty(t, out)
	require.EqualError(t, err, `required flag(s) "name", "type" not set`)

	lbName := withSuffix("test-load-balancer")
	lbID, err := createLoadBalancer(t, lbName, "--type", TestLoadBalancerTypeName, "--location", TestLocationName)
	require.NoError(t, err)

	t.Run("labels", func(t *testing.T) {
		t.Run("add-label", func(t *testing.T) {
			t.Run("non-existing", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "add-label", "non-existing-load-balancer", "foo=bar")
				require.EqualError(t, err, "load balancer not found: non-existing-load-balancer")
				assert.Empty(t, out)
			})

			t.Run("1", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "add-label", strconv.FormatInt(lbID, 10), "foo=bar")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Label(s) foo added to Load Balancer %d\n", lbID), out)
			})

			t.Run("2", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "add-label", strconv.FormatInt(lbID, 10), "baz=qux")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Label(s) baz added to Load Balancer %d\n", lbID), out)
			})
		})

		t.Run("remove-label", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "remove-label", strconv.FormatInt(lbID, 10), "baz")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Load Balancer %d\n", lbID), out)
		})
	})

	t.Run("enable-protection", func(t *testing.T) {
		t.Run("non-existing-protection", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "enable-protection", strconv.FormatInt(lbID, 10), "non-existing-protection")
			require.EqualError(t, err, "unknown protection level: non-existing-protection")
			assert.Empty(t, out)
		})

		t.Run("non-existing-load-balancer", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "enable-protection", "non-existing-load-balancer", "delete")
			require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
			assert.Empty(t, out)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "enable-protection", strconv.FormatInt(lbID, 10), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection enabled for Load Balancer %d\n", lbID), out)
		})
	})

	t.Run("change-type", func(t *testing.T) {
		t.Run("non-existing-load-balancer", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "change-type", "non-existing-load-balancer", TestLoadBalancerTypeName)
			require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "change-type", strconv.FormatInt(lbID, 10), TestLoadBalancerTypeName)
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("LoadBalancer %d changed to type %s\n", lbID, TestLoadBalancerTypeName), out)
		})
	})

	t.Run("change-algorithm", func(t *testing.T) {
		t.Run("non-existing-load-balancer", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "change-algorithm", "non-existing-load-balancer", "--algorithm-type", "round_robin")
			require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
			assert.Empty(t, out)
		})

		t.Run("non-existing-algorithm", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "change-algorithm", strconv.FormatInt(lbID, 10), "--algorithm-type", "non-existing-algorithm")
			assert.Regexp(t, `^invalid input in field 'type' \(invalid_input, [0-9a-f]+\)$`, err.Error())
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "change-algorithm", strconv.FormatInt(lbID, 10), "--algorithm-type", "round_robin")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Algorithm for Load Balancer %d was changed\n", lbID), out)
		})
	})

	t.Run("enable-public-interface", func(t *testing.T) {
		t.Run("non-existing-load-balancer", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "enable-public-interface", "non-existing-load-balancer")
			require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "enable-public-interface", strconv.FormatInt(lbID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Public interface of Load Balancer %d was enabled\n", lbID), out)
		})
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "load-balancer", "describe", strconv.FormatInt(lbID, 10))
		require.NoError(t, err)
		assert.Regexp(t, NewRegex().Start().
			Lit("ID:").Whitespace().Int().Newline().
			Lit("Name:").Whitespace().Lit(lbName).Newline().
			Lit("Created:").Whitespace().UnixDate().Lit(" (").HumanizeTime().Lit(")").Newline().
			Lit("Public Net:").Newline().
			Lit("  Enabled:").Whitespace().Lit("yes").Newline().
			Lit("  IPv4:").Whitespace().IPv4().Newline().
			Lit("  IPv4 DNS PTR:").Whitespace().AnyString().Newline().
			Lit("  IPv6:").Whitespace().IPv6().Newline().
			Lit("  IPv6 DNS PTR:").Whitespace().Lit("").Newline().
			Lit("Private Net:").Newline().
			Lit("    No Private Network").Newline().
			Lit("Algorithm:").Whitespace().Lit("round_robin").Newline().
			Lit("Load Balancer Type:").Whitespace().Lit("lb11").Lit(" (ID: 1)").Newline().
			Lit("  ID:").Whitespace().Lit("1").Newline().
			Lit("  Name:").Whitespace().Lit("lb11").Newline().
			Lit("  Description:").Whitespace().Lit("LB11").Newline().
			Lit("  Max Services:").Whitespace().Lit("5").Newline().
			Lit("  Max Connections:").Whitespace().Lit("10000").Newline().
			Lit("  Max Targets:").Whitespace().Lit("25").Newline().
			Lit("  Max assigned Certificates:").Whitespace().Lit("10").Newline().
			Lit("Services:").Newline().
			Lit("  No services").Newline().
			Lit("Targets:").Newline().
			Lit("  No targets").Newline().
			Lit("Traffic:").Newline().
			Lit("  Outgoing:").Whitespace().Lit("0 B").Newline().
			Lit("  Ingoing:").Whitespace().Lit("0 B").Newline().
			Lit("  Included:").Whitespace().Lit("20 TiB").Newline().
			Lit("Protection:").Newline().
			Lit("  Delete:").Whitespace().Lit("yes").Newline().
			Lit("Labels:").Newline().
			Lit("  foo:").Whitespace().Lit("bar").Newline().
			End(),
			out,
		)
	})

	t.Run("update-name", func(t *testing.T) {
		lbName = withSuffix("new-test-load-balancer")
		out, err := runCommand(t, "load-balancer", "update", strconv.FormatInt(lbID, 10), "--name", lbName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Load Balancer %d updated\n", lbID), out)
	})

	t.Run("disable-protection", func(t *testing.T) {
		t.Run("non-existing-load-balancer", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "disable-protection", "non-existing-load-balancer", "delete")
			require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
			assert.Empty(t, out)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "disable-protection", strconv.FormatInt(lbID, 10), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection disabled for Load Balancer %d\n", lbID), out)
		})
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "load-balancer", "delete", strconv.FormatInt(lbID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Load Balancer %d deleted\n", lbID), out)
	})
}

func createLoadBalancer(t *testing.T, name string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.LoadBalancer.Delete(context.Background(), &hcloud.LoadBalancer{Name: name})
	})

	out, err := runCommand(t, append([]string{"load-balancer", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	firstLine := strings.Split(out, "\n")[0]
	if !assert.Regexp(t, `^Load Balancer [0-9]+ created$`, firstLine) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(firstLine[14:len(firstLine)-8], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.LoadBalancer.Delete(context.Background(), &hcloud.LoadBalancer{ID: id})
	})
	return id, nil
}
