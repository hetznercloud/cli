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
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/sshutil"
)

func TestSSHKey(t *testing.T) {
	t.Parallel()

	pubKey, fingerprint, err := generateSSHKey()
	require.NoError(t, err)

	sshKeyName := withSuffix("test-ssh-key")
	sshKeyID, err := createSSHKey(t, sshKeyName, "--public-key", pubKey)
	require.NoError(t, err)

	t.Run("add-label", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "ssh-key", "add-label", "non-existing-ssh-key", "foo=bar")
			require.EqualError(t, err, "ssh key not found: non-existing-ssh-key")
			assert.Empty(t, out)
		})

		t.Run("1", func(t *testing.T) {
			out, err := runCommand(t, "ssh-key", "add-label", strconv.FormatInt(sshKeyID, 10), "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to SSH Key %d\n", sshKeyID), out)
		})

		t.Run("2", func(t *testing.T) {
			out, err := runCommand(t, "ssh-key", "add-label", strconv.FormatInt(sshKeyID, 10), "baz=qux")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz added to SSH Key %d\n", sshKeyID), out)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "ssh-key", "list", "-o=columns=id,name,fingerprint,public_key,labels,created,age")
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "NAME", "FINGERPRINT", "PUBLIC KEY", "LABELS", "CREATED", "AGE").Newline().
					Lit(strconv.FormatInt(sshKeyID, 10)).Whitespace().
					Lit(sshKeyName).Whitespace().
					Lit(fingerprint).Whitespace().
					Lit(pubKey).Whitespace().
					Lit("baz=qux, foo=bar").Whitespace().
					UnixDate().Whitespace().
					Age().Newline().
					End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			out, err := runCommand(t, "ssh-key", "list", "-o=json")
			require.NoError(t, err)
			assertjson.Equal(t, []byte(fmt.Sprintf(`
[
  {
    "id": %d,
    "name": %q,
    "fingerprint": %q,
    "public_key": %q,
    "labels": {
      "baz": "qux",
      "foo": "bar"
    },
    "created": "<ignore-diff>"
  }
]`, sshKeyID, sshKeyName, fingerprint, pubKey)), []byte(out))
		})
	})

	t.Run("remove-label", func(t *testing.T) {
		out, err := runCommand(t, "ssh-key", "remove-label", strconv.FormatInt(sshKeyID, 10), "baz")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Label(s) baz removed from SSH Key %d\n", sshKeyID), out)
	})

	t.Run("update-name", func(t *testing.T) {
		sshKeyName = withSuffix("new-test-ssh-key")
		out, err := runCommand(t, "ssh-key", "update", strconv.FormatInt(sshKeyID, 10), "--name", sshKeyName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("SSHKey %d updated\n", sshKeyID), out)
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "ssh-key", "describe", strconv.FormatInt(sshKeyID, 10))
		require.NoError(t, err)
		assert.Regexp(t, NewRegex().Start().
			Lit("ID:").Whitespace().Lit(strconv.FormatInt(sshKeyID, 10)).Newline().
			Lit("Name:").Whitespace().Lit(sshKeyName).Newline().
			Lit("Created:").Whitespace().UnixDate().Lit(" (").HumanizeTime().Lit(")").Newline().
			Lit("Fingerprint:").Whitespace().Lit(fingerprint).Newline().
			Lit("Public Key:").Newline().Lit(pubKey).
			Lit("Labels:").Newline().
			Lit("  foo:").Whitespace().Lit("bar").Newline().
			End(),
			out,
		)
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "ssh-key", "delete", strconv.FormatInt(sshKeyID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("SSH Key %d deleted\n", sshKeyID), out)
	})
}

func createSSHKey(t *testing.T, name string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.SSHKey.Delete(context.Background(), &hcloud.SSHKey{Name: name})
	})

	out, err := runCommand(t, append([]string{"ssh-key", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^SSH key [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[8:len(out)-9], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.SSHKey.Delete(context.Background(), &hcloud.SSHKey{ID: id})
	})
	return id, nil
}

func generateSSHKey() (string, string, error) {
	// ed25519 SSH key
	_, pub, err := sshutil.GenerateKeyPair()
	if err != nil {
		return "", "", err
	}

	// MD5 fingerprint
	fingerprint, err := sshutil.GetPublicKeyFingerprint(pub)
	if err != nil {
		return "", "", err
	}

	return string(pub), fingerprint, nil
}
