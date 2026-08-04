package state

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/state/config"
)

func TestDebugFileIsPrivateAndStructured(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hcloud-debug.log")
	cfg := config.New()
	require.NoError(t, cfg.LoadReader(bytes.NewReader(nil)))
	config.OptionDebug.Override(cfg, true)
	config.OptionDebugFile.Override(cfg, path)

	s, err := New(cfg, Options{Stderr: io.Discard})
	require.NoError(t, err)
	s.Logger().DebugContext(t.Context(), "diagnostic event", "component", "test")
	require.NoError(t, s.Close())

	info, err := os.Stat(path)
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
	}

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(contents), "hcloud debug session")
	assert.Contains(t, string(contents), `"msg":"diagnostic event"`)
	assert.NotContains(t, string(contents), "HCLOUD_TOKEN")
	assert.NotContains(t, string(contents), "--token")
}

func TestNewClosesDebugFileOnLaterConfigurationError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hcloud-debug.log")
	cfg := config.New()
	require.NoError(t, cfg.LoadReader(strings.NewReader("")))
	config.OptionDebug.Override(cfg, true)
	config.OptionDebugFile.Override(cfg, path)
	config.OptionPollInterval.OverrideAny(cfg, "not-a-duration")

	_, err := New(cfg, Options{Stderr: io.Discard})
	require.Error(t, err)
	require.NoError(t, os.Remove(path), "debug file must be closed when construction fails")
}
