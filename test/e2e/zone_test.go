//go:build e2e

package e2e

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestZone(t *testing.T) {
	t.Parallel()

	zoneName := withSuffix("hcloud-cli-test-zone") + ".com"

	err := createZone(t, zoneName)
	require.NoError(t, err)

	t.Run("add-label", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "zone", "add-label", "example.com", "foo=bar")
			require.EqualError(t, err, "Zone not found: example.com")
			assert.Empty(t, out)
		})

		t.Run("1", func(t *testing.T) {
			out, err := runCommand(t, "zone", "add-label", zoneName, "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to Zone %s\n", zoneName), out)
		})

		t.Run("2", func(t *testing.T) {
			out, err := runCommand(t, "zone", "add-label", zoneName, "baz=qux")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz added to Zone %s\n", zoneName), out)
		})
	})

	t.Run("remove-label", func(t *testing.T) {
		out, err := runCommand(t, "zone", "remove-label", zoneName, "baz")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Zone %s\n", zoneName), out)
	})

	t.Run("enable-protection", func(t *testing.T) {
		t.Run("unknown-protection-level", func(t *testing.T) {
			out, err := runCommand(t, "zone", "enable-protection", zoneName, "unknown-protection-level")
			require.EqualError(t, err, "unknown protection level: unknown-protection-level")
			assert.Empty(t, out)
		})

		t.Run("non-existing-zone", func(t *testing.T) {
			out, err := runCommand(t, "zone", "enable-protection", "example.com", "delete")
			require.EqualError(t, err, "Zone not found: example.com")
			assert.Empty(t, out)
		})

		t.Run("enable-delete-protection", func(t *testing.T) {
			out, err := runCommand(t, "zone", "enable-protection", zoneName, "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection enabled for Zone %s\n", zoneName), out)
		})
	})

	t.Run("change-ttl", func(t *testing.T) {
		out, err := runCommand(t, "zone", "change-ttl", zoneName, "--ttl", "600")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Changed default TTL on Zone %s\n", zoneName), out)
	})

	t.Run("add-record", func(t *testing.T) {
		out, err := runCommand(t, "zone", "add-records", zoneName, "www", "A", "--record", "192.168.0.1")
		require.NoError(t, err)
		assert.Equal(t, "Added records on Zone RRSet www A\n", out)
	})

	var zonefile string

	t.Run("export-zonefile", func(t *testing.T) {
		out, err := runCommand(t, "zone", "export-zonefile", zoneName)
		require.NoError(t, err)
		require.Contains(t, out, fmt.Sprintf("$ORIGIN\t%s.\n", zoneName))
		zonefile = out
	})

	zonefile += "\nfoobar\tIN\tA\t192.168.0.2\n"

	t.Run("import-zonefile", func(t *testing.T) {
		zonefilePath := path.Join(t.TempDir(), "example.zone")
		err := os.WriteFile(zonefilePath, []byte(zonefile), 0400)
		require.NoError(t, err)

		out, err := runCommand(t, "zone", "import-zonefile", zoneName, "--zonefile", zonefilePath)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Zone file for Zone %s imported\n", zoneName), out)
	})

	t.Run("rrset", func(t *testing.T) {
		t.Run("change-ttl", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "change-ttl", zoneName, "foobar", "A", "--ttl", "1337")
			require.NoError(t, err)
			assert.Equal(t, "Changed TTL on Zone RRSet foobar A\n", out)
		})

		t.Run("enable-protection", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "enable-protection", zoneName, "foobar", "A", "change")
			require.NoError(t, err)
			assert.Equal(t, "Resource protection enabled for Zone RRSet foobar A\n", out)
		})

		t.Run("describe", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "describe", zoneName, "foobar", "A")
			require.NoError(t, err)

			assert.Regexp(t,
				NewRegex().Start().
					Lit("ID:").Whitespace().Lit("foobar/A").Newline().
					Lit("Type:").Whitespace().Lit("A").Newline().
					Lit("Name:").Whitespace().Lit("foobar").Newline().
					Lit("TTL:").Whitespace().Lit("1337").Newline().
					Newline().
					Lit("Protection:").Newline().
					Lit("  Change:").Whitespace().Lit("yes").Newline().
					Newline().
					Lit("Labels:").Newline().
					Lit("  No labels").Newline().
					Newline().
					Lit("Records:").Newline().
					Lit("  - Value:").Whitespace().Lit("192.168.0.2").Newline().
					End(),
				out,
			)
		})

		t.Run("add-label", func(t *testing.T) {
			t.Run("non-existing", func(t *testing.T) {
				out, err := runCommand(t, "zone", "rrset", "add-label", "example.com", "foobar", "A", "foo=bar")
				require.EqualError(t, err, "Zone not found: example.com")
				assert.Empty(t, out)
			})

			t.Run("1", func(t *testing.T) {
				out, err := runCommand(t, "zone", "rrset", "add-label", zoneName, "foobar", "A", "foo=bar")
				require.NoError(t, err)
				assert.Equal(t, "Label(s) foo added to Zone RRSet foobar A\n", out)
			})

			t.Run("2", func(t *testing.T) {
				out, err := runCommand(t, "zone", "rrset", "add-label", zoneName, "foobar", "A", "baz=qux")
				require.NoError(t, err)
				assert.Equal(t, "Label(s) baz added to Zone RRSet foobar A\n", out)
			})
		})

		t.Run("remove-label", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "remove-label", zoneName, "foobar", "A", "baz")
			require.NoError(t, err)
			assert.Equal(t, "Label(s) baz removed from Zone RRSet foobar A\n", out)
		})

		t.Run("list", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "list", zoneName, "-o=columns=name,type,ttl,protection,records,labels", "--sort=name:asc")
			require.NoError(t, err)

			// TODO: Output looks good, but test regex needs to be fixed after release.
			t.SkipNow()

			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("NAME", "TYPE", "TTL", "PROTECTION", "RECORDS", "LABELS").Whitespace().Newline().
					Lit("@").Whitespace().Lit("NS").Whitespace().Lit("-").Whitespace().AnyString().Newline().
					Lit("@").Whitespace().Lit("SOA").Whitespace().Lit("-").Whitespace().AnyString().Newline().
					Lit("foobar").Whitespace().Lit("A").Whitespace().Lit("1337").Whitespace().Lit("change").
					Whitespace().Lit("192.168.0.2").Whitespace().Lit("foo=bar").Newline().
					Lit("www").Whitespace().Lit("A").Whitespace().Lit("-").Whitespace().Lit("192.168.0.1").Whitespace().Newline().
					End(),
				out,
			)
		})

		t.Run("delete-protected", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "delete", zoneName, "foobar", "A")
			assert.Regexp(t, `^RRSet\(s\) is/are change protected \(protected, [0-9a-f]+\)$`, err.Error())
			assert.Empty(t, out)
		})

		t.Run("disable-protection", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "disable-protection", zoneName, "foobar", "A", "change")
			require.NoError(t, err)
			assert.Equal(t, "Resource protection disabled for Zone RRSet foobar A\n", out)
		})

		t.Run("delete", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "delete", zoneName, "foobar", "A")
			require.NoError(t, err)
			assert.Equal(t, "Zone RRSet foobar A deleted\n", out)
		})

		t.Run("create", func(t *testing.T) {
			out, err := runCommand(t, "zone", "rrset", "create", zoneName, "--name", "foo", "--type", "A", "--record", "10.0.123.0")
			require.NoError(t, err)
			assert.Equal(t, "Zone RRSet foo A created\n", out)
		})
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "zone", "describe", zoneName)
		require.NoError(t, err)

		assert.Regexp(t,
			NewRegex().Start().
				Lit("ID:").Whitespace().Int().Newline().
				Lit("Name:").Whitespace().Lit(zoneName).Newline().
				Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
				Lit("Mode:").Whitespace().Lit("primary").Newline().
				Lit("Status:").Whitespace().Lit("ok").Newline().
				Lit("TTL:").Whitespace().Lit("600").Newline().
				Lit("Registrar:").Whitespace().Lit("other").Newline().
				Lit("Record Count:").Whitespace().Lit("6").Newline().
				Newline().
				Lit("Protection:").Newline().
				Lit("  Delete:").Whitespace().Lit("yes").Newline().
				Newline().
				Lit("Labels:").Newline().
				Lit("  foo:").Whitespace().Lit("bar").Newline().
				Newline().
				Lit("Authoritative Nameservers:").Newline().
				Lit("  Assigned:").Newline().
				AnyTimes(
					NewRegex().Lit("    - ").AnyString().Newline(),
				).
				Lit("  Delegated:").Newline().
				Lit("    No delegated nameservers").Newline().
				Lit("  Delegation last check:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
				Lit("  Delegation status:").Whitespace().OneOfLit("invalid", "unknown", "unregistered").Newline().
				End(),
			out,
		)
	})

	t.Run("disable-protection", func(t *testing.T) {
		t.Run("non-existing-zone", func(t *testing.T) {
			out, err := runCommand(t, "zone", "disable-protection", "example.com", "delete")
			require.EqualError(t, err, "Zone not found: example.com")
			assert.Empty(t, out)
		})

		t.Run("unknown-protection-level", func(t *testing.T) {
			out, err := runCommand(t, "zone", "disable-protection", zoneName, "unknown-protection-level")
			require.EqualError(t, err, "unknown protection level: unknown-protection-level")
			assert.Empty(t, out)
		})

		t.Run("disable-delete-protection", func(t *testing.T) {
			out, err := runCommand(t, "zone", "disable-protection", zoneName, "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection disabled for Zone %s\n", zoneName), out)
		})
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "zone", "delete", zoneName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Zone %s deleted\n", zoneName), out)
	})
}

func createZone(t *testing.T, name string, args ...string) error {
	t.Helper()
	t.Cleanup(func() {
		_, _, _ = client.Zone.Delete(context.Background(), &hcloud.Zone{Name: name})
	})

	_, err := runCommand(t, append([]string{"zone", "create", "--name", name}, args...)...)
	return err
}
