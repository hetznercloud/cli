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
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func TestTokenCommandProcess(_ *testing.T) {
	if os.Getenv("HCLOUD_TEST_TOKEN_COMMAND") != "1" {
		return
	}
	f, err := os.OpenFile(os.Getenv("HCLOUD_TEST_COMMAND_CALLS"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600) //nolint:gosec // The parent test supplies a path inside t.TempDir().
	if err != nil {
		os.Exit(2)
	}
	_, _ = f.WriteString("run\n")
	_ = f.Close()
	if os.Getenv("HCLOUD_TEST_COMMAND_FAIL") == "1" {
		os.Exit(1)
	}
	switch mode := os.Args[len(os.Args)-1]; mode {
	case "empty":
		fmt.Print(" \t\r\n")
	case "failure":
		fmt.Print("secret output from failed command")
		os.Exit(7)
	case "wait":
		time.Sleep(time.Minute)
	case "stdin":
		fmt.Fprint(os.Stderr, "Passphrase: ")
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			os.Exit(3)
		}
		fmt.Print(string(input))
	case "oversized":
		fmt.Print(strings.Repeat("x", 2<<20))
	case "invalid":
		fmt.Print("token\nextra line")
	default:
		fmt.Print(" \tcommand-secret\r\n")
	}
	os.Exit(0)
}

func commandState(t *testing.T) (state.State, string) {
	t.Helper()
	t.Setenv("HCLOUD_CONTEXT", "")
	t.Setenv("HCLOUD_TEST_TOKEN_COMMAND", "1")
	calls := filepath.Join(t.TempDir(), "calls")
	t.Setenv("HCLOUD_TEST_COMMAND_CALLS", calls)
	ctx := config.NewTokenCommandContext("command", testCommand("success"))
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
	require.True(t, os.IsNotExist(err), "state initialization must not execute the token command")
	return s, calls
}

func TestTokenCommandAuthentication(t *testing.T) {
	for _, ensure := range []bool{false, true} {
		t.Run(fmt.Sprint(ensure), func(t *testing.T) {
			t.Setenv("HCLOUD_TOKEN", "")
			s, calls := commandState(t)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "Bearer command-secret", r.Header.Get("Authorization"))
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
			assert.Equal(t, "run\n", string(data))
			var out bytes.Buffer
			require.NoError(t, s.Config().Write(&out))
			assert.NotContains(t, out.String(), "command-secret")
			token, err := config.OptionToken.Get(s.Config())
			require.NoError(t, err)
			assert.Empty(t, token)
		})
	}
}

func TestTokenCommandEnvironmentOverride(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "environment-token")
	t.Setenv("HCLOUD_TEST_COMMAND_FAIL", "1")
	s, calls := commandState(t)
	require.NoError(t, s.EnsureToken(nil, nil))
	req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)
	assert.Equal(t, "Bearer environment-token", req.Header.Get("Authorization"))
	_, err = os.Stat(calls)
	assert.True(t, os.IsNotExist(err))
}

func TestTokenCommandFailure(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	t.Setenv("HCLOUD_TEST_COMMAND_FAIL", "1")
	s, calls := commandState(t)
	require.ErrorContains(t, s.EnsureToken(nil, nil), "token command failed")
	require.ErrorContains(t, s.EnsureToken(nil, nil), "token command failed")
	data, err := os.ReadFile(calls)
	require.NoError(t, err)
	assert.Equal(t, "run\n", string(data))
}

func TestManagementDoesNotRetrieveCredentials(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"version"}, {"context", "list", "-o", "json", "--allow-sensitive"}, {"config", "get", "token", "--allow-sensitive"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			t.Setenv("HCLOUD_TOKEN", "")
			t.Setenv("HCLOUD_TEST_COMMAND_FAIL", "1")
			s, calls := commandState(t)
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

func TestTokenCommandRedirectDoesNotLeakCredentials(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	s, _ := commandState(t)
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, "{}")
	}))
	defer destination.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer command-secret", r.Header.Get("Authorization"))
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

func TestTokenCommandFailurePreventsAPIRequest(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	t.Setenv("HCLOUD_TEST_COMMAND_FAIL", "1")
	s, _ := commandState(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Error("failed token command must prevent the API request")
		_, _ = io.WriteString(w, "{}")
	}))
	defer server.Close()
	req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)
	target, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.URL = target.URL
	_, err = s.Client().Do(req, io.Discard)
	require.ErrorContains(t, err, "token command failed")
}

// Execute this test binary as the secret producer, with no credential protocol.
func testCommand(mode string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(`"%s" -test.run=^TestTokenCommandProcess$ -- %s`, os.Args[0], mode)
	}
	return "'" + strings.ReplaceAll(os.Args[0], "'", "'\\''") + "' -test.run=^TestTokenCommandProcess$ -- " + mode
}

func TestTokenCommandOutputFailures(t *testing.T) {
	for _, mode := range []string{"empty", "failure", "oversized", "invalid"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("HCLOUD_TOKEN", "")
			s, _ := commandState(t)
			cfg := s.Config()
			config.OptionToken.Override(cfg, "stored-token")
			ctx := config.NewTokenCommandContext("test", testCommand(mode))
			cfg.SetActiveContext(ctx)
			// A new client lifecycle uses the selected command.
			s, err := state.New(context.Background(), cfg)
			require.NoError(t, err)
			err = s.EnsureToken(nil, nil)
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "secret output")
			if mode == "failure" {
				assert.Contains(t, err.Error(), "exit code 7")
			}
			if mode == "empty" {
				assert.Contains(t, err.Error(), "empty")
			}
		})
	}
}

func TestTokenCommandDebugDoesNotExposeToken(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	s, _ := commandState(t)
	cfg := s.Config()
	debugFile := filepath.Join(t.TempDir(), "debug.log")
	config.OptionDebug.Override(cfg, true)
	config.OptionDebugFile.Override(cfg, debugFile)
	s, err := state.New(context.Background(), cfg)
	require.NoError(t, err)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer command-secret", r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, "{}")
	}))
	defer server.Close()
	req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)
	target, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.URL = target.URL
	_, err = s.Client().Do(req, io.Discard)
	require.NoError(t, err)
	output, err := os.ReadFile(debugFile)
	require.NoError(t, err)
	assert.NotContains(t, string(output), "command-secret")
	assert.Contains(t, string(output), "REDACTED")
}

func TestTokenCommandCancellation(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	s, calls := commandState(t)
	cfg := s.Config()
	cfg.SetActiveContext(config.NewTokenCommandContext("test", testCommand("wait")))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s, err := state.New(ctx, cfg)
	require.NoError(t, err)
	require.ErrorIs(t, s.EnsureToken(nil, nil), context.Canceled)
	_, err = os.Stat(calls)
	assert.True(t, os.IsNotExist(err))
}

func TestTokenCommandShellComposition(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell syntax")
	}
	t.Setenv("HCLOUD_TOKEN", "")
	s, _ := commandState(t)
	cfg := s.Config()
	cfg.SetActiveContext(config.NewTokenCommandContext("test", `printf '%s' "$HCLOUD_TEST_PIPE_TOKEN" | tr '[:lower:]' '[:upper:]'`))
	t.Setenv("HCLOUD_TEST_PIPE_TOKEN", "pipe-token")
	s, err := state.New(context.Background(), cfg)
	require.NoError(t, err)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer PIPE-TOKEN", r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, "{}")
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

func TestTokenCommandInheritsIO(t *testing.T) {
	t.Setenv("HCLOUD_TOKEN", "")
	s, _ := commandState(t)
	dir := t.TempDir()
	input, err := os.CreateTemp(dir, "stdin")
	require.NoError(t, err)
	defer input.Close()
	_, err = input.WriteString(" \tcommand-secret\r\n")
	require.NoError(t, err)
	_, err = input.Seek(0, io.SeekStart)
	require.NoError(t, err)
	diagnostic, err := os.CreateTemp(dir, "stderr")
	require.NoError(t, err)
	defer diagnostic.Close()
	originalStdin, originalStderr := os.Stdin, os.Stderr
	os.Stdin, os.Stderr = input, diagnostic
	defer func() { os.Stdin, os.Stderr = originalStdin, originalStderr }()
	cfg := s.Config()
	cfg.SetActiveContext(config.NewTokenCommandContext("test", testCommand("stdin")))
	s, err = state.New(context.Background(), cfg)
	require.NoError(t, err)
	require.NoError(t, s.EnsureToken(nil, nil))
	output, err := os.ReadFile(diagnostic.Name())
	require.NoError(t, err)
	assert.Equal(t, "Passphrase: ", string(output))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer command-secret", r.Header.Get("Authorization"))
		_, _ = io.WriteString(w, "{}")
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

func TestTokenCommandPrecedence(t *testing.T) {
	for _, tt := range []struct{ command, expected string }{
		{"", "stored-token"},
		{" \t", "stored-token"},
		{testCommand("success"), "command-secret"},
	} {
		t.Run(tt.command, func(t *testing.T) {
			t.Setenv("HCLOUD_TOKEN", "")
			_, _ = commandState(t)
			cfg := config.New()
			require.NoError(t, cfg.Read([]byte(fmt.Sprintf("active_context = \"test\"\n[[contexts]]\nname = \"test\"\ntoken = \"stored-token\"\ntoken_command = %q\n", tt.command))))
			s, err := state.New(context.Background(), cfg)
			require.NoError(t, err)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "Bearer "+tt.expected, r.Header.Get("Authorization"))
				_, _ = io.WriteString(w, "{}")
			}))
			defer server.Close()
			req, err := s.Client().NewRequest(context.Background(), http.MethodGet, "", nil)
			require.NoError(t, err)
			target, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
			require.NoError(t, err)
			req.URL = target.URL
			_, err = s.Client().Do(req, io.Discard)
			require.NoError(t, err)
		})
	}
}
