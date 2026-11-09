package config

import (
	"bytes"
	stdcontext "context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Run the test binary itself as a portable credential helper.
func TestCredentialHelperProcess(_ *testing.T) {
	if os.Getenv("HCLOUD_TEST_CREDENTIAL_HELPER") != "1" {
		return
	}
	args := os.Args
	mode := args[len(args)-2]
	if args[len(args)-1] != "get" {
		os.Exit(2)
	}
	input, _ := io.ReadAll(os.Stdin)
	if string(input) != "protocol=https\nhost=api.hetzner.cloud\npath=production\nusername=alice\n\n" {
		os.Exit(3)
	}
	switch mode {
	case "success":
		fmt.Print("username=alice\nunknown=ignored\npassword=secret=token\n\n")
	case "empty":
		fmt.Print("username=alice\n\n")
	case "malformed":
		fmt.Print("secret token without framing\n")
	case "oversized":
		fmt.Print(strings.Repeat("x", 2<<20))
	case "crlf":
		fmt.Print("password=token\r\n\r\n")
	case "failure":
		fmt.Fprint(os.Stderr, "do not expose this secret")
		os.Exit(1)
	case "wait":
		time.Sleep(time.Minute)
	}
	os.Exit(0)
}

func TestResolveCredential(t *testing.T) {
	t.Setenv("HCLOUD_TEST_CREDENTIAL_HELPER", "1")
	for _, mode := range []string{"success", "empty", "malformed", "failure", "wait", "oversized", "crlf"} {
		t.Run(mode, func(t *testing.T) {
			ctx := stdcontext.Background()
			if mode == "wait" {
				var cancel stdcontext.CancelFunc
				ctx, cancel = stdcontext.WithTimeout(ctx, 100*time.Millisecond)
				defer cancel()
			}
			token, err := ResolveCredential(ctx, []string{os.Args[0], "-test.run=^TestCredentialHelperProcess$", "--", mode}, "https://alice@api.hetzner.cloud/production")
			if mode == "success" || mode == "crlf" {
				require.NoError(t, err)
				if mode == "success" {
					assert.Equal(t, "secret=token", token)
				} else {
					assert.Equal(t, "token", token)
				}
			} else {
				require.Error(t, err)
				assert.Empty(t, token)
				assert.NotContains(t, err.Error(), "do not expose this secret")
				assert.NotContains(t, err.Error(), "secret token")
			}
		})
	}
}

func TestCredentialContextValidation(t *testing.T) {
	for _, url := range []string{"not-a-url", "https:///missing-host", "https://alice:secret@example.com", "https://example.com/?query", "https://example.com/#fragment", "https://example.com/%0apassword=bad", "https://example.com/%00bad"} {
		t.Run(url, func(t *testing.T) {
			_, err := NewCredentialContext("test", []string{"helper"}, url)
			require.Error(t, err)
		})
	}
	for _, fields := range []string{
		`credential_helper = ["helper"]`,
		`credential_url = "https://example.com"`,
		"credential_helper = [\"helper\"]\ncredential_url = \"https://example.com\"\ntoken = \"stored\"",
	} {
		cfg := New()
		require.Error(t, cfg.Read([]byte("[[contexts]]\nname = \"test\"\n"+fields)))
	}
}

func TestCredentialContextRoundTrip(t *testing.T) {
	t.Setenv("HCLOUD_CONTEXT", "")
	t.Setenv("HCLOUD_TOKEN", "")
	ctx, err := NewCredentialContext("test", []string{"missing-helper", "argument with spaces"}, "https://example.com/resource")
	require.NoError(t, err)
	cfg := New()
	cfg.SetContexts([]Context{ctx})
	cfg.SetActiveContext(ctx)
	RenameContext(ctx, "renamed")
	cfg.SetActiveContext(ctx)
	var out bytes.Buffer
	require.NoError(t, cfg.Write(&out))
	assert.NotContains(t, out.String(), "token =")
	loaded := New()
	require.NoError(t, loaded.Read(out.Bytes()))
	assert.Equal(t, "renamed", loaded.ActiveContext().Name())
	assert.Equal(t, ctx.CredentialHelper(), loaded.ActiveContext().CredentialHelper())
	assert.Equal(t, ctx.CredentialURL(), loaded.ActiveContext().CredentialURL())
	assert.Empty(t, loaded.ActiveContext().Token())
}

func TestMissingCredentialHelper(t *testing.T) {
	_, err := ResolveCredential(stdcontext.Background(), []string{t.TempDir() + "/missing-helper"}, "https://example.com/resource")
	require.EqualError(t, err, "unable to execute credential helper")
}
