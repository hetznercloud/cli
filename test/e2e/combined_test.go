//go:build e2e

package e2e

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/sshutil"
)

func TestCombined(t *testing.T) {
	// combined tests combine multiple resources and can thus not be run in parallel
	priv, pub, err := sshutil.GenerateKeyPair()
	require.NoError(t, err)

	keyDir := t.TempDir()
	pubKeyPath, privKeyPath := path.Join(keyDir, "id_ed25519.pub"), path.Join(keyDir, "id_ed25519")
	err = os.WriteFile(privKeyPath, priv, 0600)
	require.NoError(t, err)
	err = os.WriteFile(pubKeyPath, pub, 0644)
	require.NoError(t, err)

	networkName := withSuffix("test-network")
	networkID, err := createNetwork(t, networkName, "--ip-range", "10.0.0.0/16")
	require.NoError(t, err)

	t.Run("network", func(t *testing.T) {
		t.Run("add-subnet", func(t *testing.T) {
			out, err := runCommand(t, "network", "add-subnet", "--type", "cloud", "--network-zone", TestNetworkZone, "--ip-range", "10.0.1.0/24", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Subnet added to network %d\n", networkID), out)
		})
	})

	sshKeyName := withSuffix("test-ssh-key")
	sshKeyID, err := createSSHKey(t, sshKeyName, "--public-key-from-file", pubKeyPath)
	require.NoError(t, err)

	serverName := withSuffix("test-server")
	serverID, err := createServer(t, serverName, TestServerType, TestImage, "--ssh-key", strconv.FormatInt(sshKeyID, 10), "--network", strconv.FormatInt(networkID, 10))
	require.NoError(t, err)

	firewallName := withSuffix("test-firewall")
	firewallID, err := createFirewall(t, firewallName)
	require.NoError(t, err)

	loadBalancerName := withSuffix("test-load-balancer")
	loadBalancerID, err := createLoadBalancer(t, loadBalancerName, "--location", TestLocationName, "--type", TestLoadBalancerTypeName, "--network", strconv.FormatInt(networkID, 10))
	require.NoError(t, err)

	t.Run("load-balancer", func(t *testing.T) {
		t.Run("detach-from-network", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "detach-from-network", strconv.FormatInt(loadBalancerID, 10), "--network", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Load Balancer %d detached from Network %d\n", loadBalancerID, networkID), out)
		})

		t.Run("attach-to-network", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "attach-to-network", strconv.FormatInt(loadBalancerID, 10), "--network", strconv.FormatInt(networkID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Load Balancer %d attached to network %d\n", loadBalancerID, networkID), out)
		})

		t.Run("disable-public-interface", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "disable-public-interface", strconv.FormatInt(loadBalancerID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Public interface of Load Balancer %d was disabled\n", loadBalancerID), out)
		})

		t.Run("add-target", func(t *testing.T) {
			t.Run("non-existing-load-balancer", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "add-target", "non-existing-load-balancer", "--server", "my-server")
				require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
				assert.Empty(t, out)
			})

			t.Run("label-selector", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "add-target", strconv.FormatInt(loadBalancerID, 10), "--label-selector", "foo=bar")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Target added to Load Balancer %d\n", loadBalancerID), out)
			})

			t.Run("server", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "add-target", strconv.FormatInt(loadBalancerID, 10), "--server", strconv.FormatInt(serverID, 10), "--use-private-ip")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Target added to Load Balancer %d\n", loadBalancerID), out)
			})
		})

		t.Run("add-service", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "add-service", "--protocol", "http", "--listen-port", "80", "--destination-port", "80", strconv.FormatInt(loadBalancerID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Service was added to Load Balancer %d\n", loadBalancerID), out)
		})

		t.Run("update-service", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "update-service", strconv.FormatInt(loadBalancerID, 10), "--listen-port", "80", "--destination-port", "8080")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Service 80 on Load Balancer %d was updated\n", loadBalancerID), out)
		})

		t.Run("delete-service", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer", "delete-service", strconv.FormatInt(loadBalancerID, 10), "--listen-port", "80")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Service on port 80 deleted from Load Balancer %d\n", loadBalancerID), out)
		})

		t.Run("remove-target", func(t *testing.T) {
			t.Run("non-existing-load-balancer", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "remove-target", "non-existing-load-balancer", "--server", "my-server")
				require.EqualError(t, err, "Load Balancer not found: non-existing-load-balancer")
				assert.Empty(t, out)
			})

			t.Run("label-selector", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "remove-target", strconv.FormatInt(loadBalancerID, 10), "--label-selector", "foo=bar")
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Target removed from Load Balancer %d\n", loadBalancerID), out)
			})

			t.Run("server", func(t *testing.T) {
				out, err := runCommand(t, "load-balancer", "remove-target", strconv.FormatInt(loadBalancerID, 10), "--server", strconv.FormatInt(serverID, 10))
				require.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("Target removed from Load Balancer %d\n", loadBalancerID), out)
			})
		})
	})

	t.Run("firewall", func(t *testing.T) {
		t.Run("apply-to-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "apply-to-resource", "--type", "server", "--server", serverName, strconv.FormatInt(firewallID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall %d applied to resource\n", firewallID), out)
		})

		t.Run("delete-in-use", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete", strconv.FormatInt(firewallID, 10))
			assert.Regexp(t, `^firewall with ID [0-9]+ is still in use \(resource_in_use, [0-9a-f]+\)$`, err.Error())
			assert.Empty(t, out)
		})

		t.Run("remove-from-server", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "remove-from-resource", "--type", "server", "--server", serverName, strconv.FormatInt(firewallID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall %d removed from resource\n", firewallID), out)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "firewall", "delete", strconv.FormatInt(firewallID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Firewall %d deleted\n", firewallID), out)
		})

	})

	floatingIPName := withSuffix("test-floating-ip")
	floatingIP, err := createFloatingIP(t, floatingIPName, "ipv4", "--server", strconv.FormatInt(serverID, 10))
	require.NoError(t, err)

	t.Run("floating-ip", func(t *testing.T) {
		t.Run("unassign", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "unassign", strconv.FormatInt(floatingIP, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Floating IP %d unassigned\n", floatingIP), out)
		})

		t.Run("assign", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "assign", strconv.FormatInt(floatingIP, 10), strconv.FormatInt(serverID, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Floating IP %d assigned to Server %d\n", floatingIP, serverID), out)
		})

		t.Run("describe", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIP, 10))
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					Lit("ID:").Whitespace().Int().Newline().
					Lit("Type:").Whitespace().Lit("ipv4").Newline().
					Lit("Name:").Whitespace().Raw(`test-floating-ip-[0-9a-f]{8}`).Newline().
					Lit("Description:").Whitespace().Lit("-").Newline().
					Lit("Created:").Whitespace().UnixDate().Lit(" (").HumanizeTime().Lit(")").Newline().
					Lit("IP:").Whitespace().IPv4().Newline().
					Lit("Blocked:").Whitespace().Lit("no").Newline().
					Lit("Home Location:").Whitespace().LocationName().Newline().
					Lit("Server:").Newline().
					Lit("  ID:").Whitespace().Int().Newline().
					Lit("  Name:").Whitespace().Raw(`test-server-[0-9a-f]{8}`).Newline().
					Lit("DNS:").Newline().
					Lit("  ").IPv4().Lit(": static.").IPv4().Lit(".clients.your-server.de").Newline().
					Lit("Protection:").Newline().
					Lit("  Delete:").Whitespace().Lit("no").Newline().
					Lit("Labels:").Newline().
					Lit("  No labels").Newline().
					End(),
				out,
			)
		})

		t.Run("list", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "list", "-o", "columns=server", "-o", "noheader")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("%s\n", serverName), out)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIP, 10))
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIP), out)
		})
	})

	t.Run("ssh", func(t *testing.T) {
		out, err := runCommand(
			t, "server", "ssh", strconv.FormatInt(serverID, 10),
			"-i", privKeyPath,
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "IdentitiesOnly=yes",
			"--", "exit",
		)
		require.NoError(t, err)
		assert.Empty(t, out)
	})

	t.Run("delete-server", func(t *testing.T) {
		out, err := runCommand(t, "server", "delete", strconv.FormatInt(serverID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Server %d deleted\n", serverID), out)
	})

	t.Run("delete-ssh-key", func(t *testing.T) {
		out, err := runCommand(t, "ssh-key", "delete", strconv.FormatInt(sshKeyID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("SSH Key %d deleted\n", sshKeyID), out)
	})

	t.Run("delete-network", func(t *testing.T) {
		out, err := runCommand(t, "network", "delete", strconv.FormatInt(networkID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Network %d deleted\n", networkID), out)
	})

	t.Run("delete-load-balancer", func(t *testing.T) {
		out, err := runCommand(t, "load-balancer", "delete", strconv.FormatInt(loadBalancerID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Load Balancer %d deleted\n", loadBalancerID), out)
	})
}
