package config_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/state/config"
)

func TestLoadReaderUsesInjectedArguments(t *testing.T) {
	cfg := config.New(config.WithArgs([]string{"--debug"}))

	require.NoError(t, cfg.LoadReader(bytes.NewReader(nil)))
	debug, err := config.OptionDebug.Get(cfg)
	require.NoError(t, err)
	assert.True(t, debug)
}

func TestCustomOptionsAreScopedToConfigInstance(t *testing.T) {
	custom, cleanup := config.NewTestOption("custom", "custom option", "", config.OptionFlagPreference, nil)
	t.Cleanup(cleanup)
	contents := bytes.NewBufferString("[preferences]\ncustom = \"value\"\n")

	withCustom := config.New(config.WithOptions(custom))
	require.NoError(t, withCustom.LoadReader(contents))
	_, ok := withCustom.LookupOption("custom")
	assert.True(t, ok)

	withoutCustom := config.New()
	err := withoutCustom.LoadReader(bytes.NewBufferString("[preferences]\ncustom = \"value\"\n"))
	require.EqualError(t, err, "unknown preference: custom")
	_, ok = withoutCustom.LookupOption("custom")
	assert.False(t, ok)
}

func TestWriteAtomicallyCreatesPrivateConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "cli.toml")
	cfg := config.New()
	require.NoError(t, cfg.LoadFile(path))
	cfg.Preferences().Set("debug", true)

	require.NoError(t, cfg.Write(nil))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(contents), "debug = true")
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(path), ".hcloud-config-*"))
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestWriteAtomicallyReplacesExistingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cli.toml")
	require.NoError(t, os.WriteFile(path, []byte("[preferences]\ndebug = false\n"), 0600))

	cfg := config.New()
	require.NoError(t, cfg.LoadFile(path))
	cfg.Preferences().Set("debug", true)
	require.NoError(t, cfg.Write(nil))

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(contents), "debug = false")
	assert.Contains(t, string(contents), "debug = true")
}

func TestLoadReaderReturnsWarningWriteError(t *testing.T) {
	warningErr := errors.New("warning writer failed")
	cfg := config.New(config.WithWarningWriter(errorWriter{err: warningErr}))

	err := cfg.LoadReader(bytes.NewBufferString("active_context = \"missing\"\n"))
	require.ErrorIs(t, err, warningErr)
}

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}

var _ io.Writer = errorWriter{}
