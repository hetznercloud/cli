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

	sshKeyName := withSuffix("test-ssh-key")
	sshKeyID, err := createSSHKey(t, sshKeyName, "--public-key-from-file", pubKeyPath)
	require.NoError(t, err)

	serverName := withSuffix("test-server")
	serverID, err := createServer(t, serverName, TestServerType, TestImage, "--ssh-key", strconv.FormatInt(sshKeyID, 10))
	require.NoError(t, err)

	firewallName := withSuffix("test-firewall")
	firewallID, err := createFirewall(t, firewallName)
	require.NoError(t, err)

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
			assert.Equal(t, fmt.Sprintf("firewall %d deleted\n", firewallID), out)
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
			assert.Equal(t, fmt.Sprintf("Floating IP %d assigned to server %d\n", floatingIP, serverID), out)
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
}
