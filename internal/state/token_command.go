package state

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/http/httpguts"
)

// resolveTokenCommand captures only stdout; interactive input and diagnostics
// remain attached to the parent's terminal. Do not include the command or its
// stdout in errors, as either could contain secrets.
func resolveTokenCommand(ctx context.Context, command string) (string, error) {
	name, args := "/bin/sh", []string{"-c", command}
	if runtime.GOOS == "windows" {
		name, args = "cmd.exe", []string{"/D", "/C", command}
	}
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // Shell execution of the user-configured token command is intentional.
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	// A shell child may retain stdout after cancellation. Bound the wait for its
	// pipe so the CLI can still return when the command's context is canceled.
	cmd.WaitDelay = time.Second
	var output tokenOutput
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("token command canceled: %w", ctx.Err())
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("token command failed (exit code %d)", exitErr.ExitCode())
		}
		return "", errors.New("unable to execute token command")
	}
	if output.exceeded {
		return "", errors.New("token command output is too large")
	}
	token := strings.TrimSpace(output.buffer.String())
	if token == "" {
		return "", errors.New("token command returned empty output")
	}
	if !httpguts.ValidHeaderFieldValue(token) {
		return "", errors.New("token command returned invalid token characters")
	}
	return token, nil
}

// Drain excessive output without retaining it or blocking the child process.
type tokenOutput struct {
	buffer   bytes.Buffer
	exceeded bool
}

func (out *tokenOutput) Write(p []byte) (int, error) {
	const limit = 1 << 20
	n := len(p)
	if remaining := limit - out.buffer.Len(); len(p) > remaining {
		out.exceeded = true
		p = p[:remaining]
	}
	_, _ = out.buffer.Write(p)
	return n, nil
}

// The SDK puts this marker on requests instead of a secret. Resolving it in the
// transport keeps client construction and configuration inspection side-effect free.
// Only replace this header when it is present: net/http removes Authorization on
// redirects to unrelated hosts, and we must not restore it in that case.
const tokenCommandMarker = "hcloud-token-command" //nolint:gosec // Non-secret marker replaced by the transport.

type tokenResolver struct {
	once    sync.Once
	command string
	token   string
	err     error
}

func (r *tokenResolver) resolve(ctx context.Context) (string, error) {
	r.once.Do(func() {
		r.token, r.err = resolveTokenCommand(ctx, r.command)
	})
	return r.token, r.err
}

type tokenTransport struct {
	resolver *tokenResolver
	base     http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Authorization") != "Bearer "+tokenCommandMarker {
		return t.base.RoundTrip(req)
	}
	token, err := t.resolver.resolve(req.Context())
	if err != nil {
		return nil, err
	}
	cloned := req.Clone(req.Context())
	cloned.Header.Set("Authorization", "Bearer "+token)
	return t.base.RoundTrip(cloned)
}
