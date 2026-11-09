package config

import (
	"bufio"
	"bytes"
	stdcontext "context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// NewCredentialContext creates a context backed by a read-only Git credential helper.
func NewCredentialContext(name string, helper []string, credentialURL string) (Context, error) {
	ctx := &context{
		ContextName:             name,
		ContextCredentialHelper: append([]string(nil), helper...),
		ContextCredentialURL:    credentialURL,
	}
	if len(helper) == 0 || credentialURL == "" {
		return nil, errors.New("credential helper and credential URL are required")
	}
	if err := ctx.validateCredentials(); err != nil {
		return nil, err
	}
	return ctx, nil
}

func (ctx *context) validateCredentials() error {
	if len(ctx.ContextCredentialHelper) == 0 && ctx.ContextCredentialURL == "" {
		return nil
	}
	if ctx.ContextToken != "" {
		return errors.New("token and credential helper settings are mutually exclusive")
	}
	if len(ctx.ContextCredentialHelper) == 0 || strings.TrimSpace(ctx.ContextCredentialHelper[0]) == "" || ctx.ContextCredentialURL == "" {
		return errors.New("credential helper and credential URL are required")
	}
	_, err := credentialRequest(ctx.ContextCredentialURL)
	return err
}

func credentialRequest(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Hostname() == "" || u.Opaque != "" || strings.ContainsAny(rawURL, "?#\r\n\x00") {
		return "", errors.New("invalid credential URL: expected a URL with a scheme and host, without a query or fragment")
	}
	if _, hasPassword := u.User.Password(); hasPassword {
		return "", errors.New("credential URL must not contain a password")
	}
	fields := [][2]string{{"protocol", u.Scheme}, {"host", u.Host}}
	if u.Path != "" {
		fields = append(fields, [2]string{"path", strings.TrimPrefix(u.Path, "/")})
	}
	if u.User != nil {
		fields = append(fields, [2]string{"username", u.User.Username()})
	}
	var request strings.Builder
	for _, field := range fields {
		if strings.ContainsAny(field[1], "\r\n\x00") {
			return "", errors.New("credential URL contains invalid field characters")
		}
		fmt.Fprintf(&request, "%s=%s\n", field[0], field[1])
	}
	request.WriteByte('\n')
	return request.String(), nil
}

// ResolveCredential retrieves a token from the password field of a Git helper response.
// Helper output is never included in errors, since either stream may contain secrets.
func ResolveCredential(ctx stdcontext.Context, helper []string, credentialURL string) (string, error) {
	if len(helper) == 0 || strings.TrimSpace(helper[0]) == "" {
		return "", errors.New("credential helper executable is required")
	}
	request, err := credentialRequest(credentialURL)
	if err != nil {
		return "", err
	}
	args := append(append([]string(nil), helper[1:]...), "get")
	cmd := exec.CommandContext(ctx, helper[0], args...) //nolint:gosec // Execute the user-configured helper directly, without a shell.
	cmd.Stdin = strings.NewReader(request)
	var output credentialOutput
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("credential helper canceled: %w", ctx.Err())
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("credential helper failed (exit code %d)", exitErr.ExitCode())
		}
		return "", errors.New("unable to execute credential helper")
	}
	if output.exceeded {
		return "", errors.New("credential helper response is too large")
	}
	var token string
	scanner := bufio.NewScanner(bytes.NewReader(output.Bytes()))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || key == "" || strings.ContainsAny(line, "\r\x00") {
			return "", errors.New("invalid credential helper response")
		}
		if key == "password" {
			token = value
		}
	}
	if scanner.Err() != nil {
		return "", errors.New("invalid credential helper response")
	}
	if strings.TrimSpace(token) == "" {
		return "", errors.New("credential helper returned no password")
	}
	return token, nil
}

// Drain excessive output without retaining it or blocking the child process.
type credentialOutput struct {
	bytes.Buffer
	exceeded bool
}

func (out *credentialOutput) Write(p []byte) (int, error) {
	const limit = 1 << 20
	n := len(p)
	if remaining := limit - out.Len(); len(p) > remaining {
		out.exceeded = true
		p = p[:remaining]
	}
	_, _ = out.Buffer.Write(p)
	return n, nil
}
