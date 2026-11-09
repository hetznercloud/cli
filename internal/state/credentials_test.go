package state_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func TestStateCredentialHelperProcess(_ *testing.T) {
	if os.Getenv("HCLOUD_TEST_STATE_HELPER") != "1" {
		return
	}
	f, err := os.OpenFile(os.Getenv("HCLOUD_TEST_HELPER_CALLS"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600) //nolint:gosec // The parent test supplies a path inside t.TempDir().
	if err != nil {
		os.Exit(2)
	}
	_, _ = f.WriteString("get\n")
	_ = f.Close()
	_, _ = io.Copy(io.Discard, os.Stdin)
	if os.Args[len(os.Args)-1] != "get" {
		os.Exit(3)
	}
	if os.Getenv("HCLOUD_TEST_HELPER_FAIL") == "1" {
		os.Exit(1)
	}
	fmt.Print("password=helper-secret\n\n")
	os.Exit(0)
}

func helperState(t *testing.T) (state.State, string) {
	t.Helper()
	t.Setenv("HCLOUD_CONTEXT", "")
	t.Setenv("HCLOUD_TEST_STATE_HELPER", "1")
	calls := filepath.Join(t.TempDir(), "calls")
	t.Setenv("HCLOUD_TEST_HELPER_CALLS", calls)
	ctx, err := config.NewCredentialContext("helper", []string{os.Args[0], "-test.run=^TestStateCredentialHelperProcess$", "--"}, "https://api.hetzner.cloud/project")
	require.NoError(t, err)
	cfg := config.New()
	cfg.SetContexts([]config.Context{ctx})
	cfg.SetActiveContext(ctx)
	var data bytes.Buffer
	require.NoError(t, cfg.Write(&data))
	cfg = config.New()
	require.NoError(t, cfg.Read(data.Bytes()))
	s, err := state.New(context.Background(), cfg)
	require.NoError(t, err)
	_, err = os.Stat(calls)
	require.True(t, os.IsNotExist(err), "state initialization must not invoke the helper")
	return s, calls
}

func TestHelperAuthentication(t *testing.T) {
	for _, ensure := range []bool{false, true} {
		t.Run(fmt.Sprint(ensure), func(t *testing.T) {
			t.Setenv("HCLOUD_TOKEN", "")
			s, calls := helperState(t)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "Bearer helper-secret", r.Header.Get("Authorization"))
				_, _ = io.WriteString(w, "{}")
			}))
			defer server.Close()
			if ensure {
				require.NoError(t, s.EnsureToken(nil, nil))
				require.NoError(t, s.EnsureToken(nil, nil))
			}
			for range 2 {
				req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
				require.NoError(t, err)
				target, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
				require.NoError(t, err)
				req.URL = target.URL
				_, err = s.Client().Do(req, io.Discard)
				require.NoError(t, err)
			}
			data, err := os.ReadFile(calls)
			require.NoError(t, err)
			assert.Equal(t, "get\n", string(data))
			var out bytes.Buffer
			require.NoError(t, s.Config().Write(&out))
			assert.NotContains(t, out.String(), "helper-secret")
			token, err := config.OptionToken.Get(s.Config())
			require.NoError(t, err)
			assert.Empty(t, token)
		})
	}
}

func TestHelperEnvironmentOverride(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "environment-token")
	t.Setenv("HCLOUD_TEST_HELPER_FAIL", "1")
	s, calls := helperState(t)
	require.NoError(t, s.EnsureToken(nil, nil))
	req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)
	assert.Equal(t, "Bearer environment-token", req.Header.Get("Authorization"))
	_, err = os.Stat(calls)
	assert.True(t, os.IsNotExist(err))
}

func TestHelperFailure(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	t.Setenv("HCLOUD_TEST_HELPER_FAIL", "1")
	s, calls := helperState(t)
	require.ErrorContains(t, s.EnsureToken(nil, nil), "credential helper failed")
	require.ErrorContains(t, s.EnsureToken(nil, nil), "credential helper failed")
	data, err := os.ReadFile(calls)
	require.NoError(t, err)
	assert.Equal(t, "get\n", string(data))
}

func TestManagementDoesNotRetrieveCredentials(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"version"}, {"context", "list", "-o", "json", "--allow-sensitive"}, {"config", "get", "token", "--allow-sensitive"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			t.Setenv("HCLOUD_TOKEN", "")
			t.Setenv("HCLOUD_TEST_HELPER_FAIL", "1")
			s, calls := helperState(t)
			cmd := cli.NewRootCommand(s)
			cmd.SetArgs(args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			require.NoError(t, cmd.Execute())
			_, err := os.Stat(calls)
			assert.True(t, os.IsNotExist(err))
		})
	}
}

func TestHelperRedirectDoesNotLeakCredentials(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	s, _ := helperState(t)
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, "{}")
	}))
	defer destination.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer helper-secret", r.Header.Get("Authorization"))
		// A different hostname makes net/http strip the Authorization header.
		http.Redirect(w, r, strings.Replace(destination.URL, "127.0.0.1", "localhost", 1), http.StatusFound)
	}))
	defer server.Close()
	req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)
	target, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.URL = target.URL
	_, err = s.Client().Do(req, io.Discard)
	require.NoError(t, err)
}

func TestHelperFailurePreventsAPIRequest(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	t.Setenv("HCLOUD_TEST_HELPER_FAIL", "1")
	s, _ := helperState(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("failed credential lookup must prevent the API request")
		_, _ = io.WriteString(w, "{}")
	}))
	defer server.Close()
	req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)
	target, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.URL = target.URL
	_, err = s.Client().Do(req, io.Discard)
	require.ErrorContains(t, err, "credential helper failed")
}
