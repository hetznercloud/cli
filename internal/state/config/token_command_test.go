package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenCommandContextRoundTrip(t *testing.T) {
	t.Setenv("HCLOUD_CONTEXT", "")
	t.Setenv("HCLOUD_TOKEN", "")
	command := `passbolt get resource --id example --field password --raw`
	ctx := NewTokenCommandContext("test", command)
	cfg := New()
	cfg.SetContexts([]Context{ctx})
	cfg.SetActiveContext(ctx)
	RenameContext(ctx, "renamed")
	cfg.SetActiveContext(ctx)
	var out bytes.Buffer
	require.NoError(t, cfg.Write(&out))
	assert.NotContains(t, out.String(), "token =")
	assert.Contains(t, out.String(), "token_command =")
	loaded := New()
	require.NoError(t, loaded.Read(out.Bytes()))
	assert.Equal(t, "renamed", loaded.ActiveContext().Name())
	assert.Equal(t, command, loaded.ActiveContext().TokenCommand())
	assert.Empty(t, loaded.ActiveContext().Token())
}

func TestTokenCommandWithStaticToken(t *testing.T) {
	t.Setenv("HCLOUD_CONTEXT", "")
	cfg := New()
	require.NoError(t, cfg.Read([]byte(`active_context = "test"
[[contexts]]
name = "test"
token = "static-token"
token_command = "printf command-token"
`)))
	assert.Equal(t, "static-token", cfg.ActiveContext().Token())
	assert.Equal(t, "printf command-token", cfg.ActiveContext().TokenCommand())
}
