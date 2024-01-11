//go:build windows

package config

import (
	"os"
	"path/filepath"
)

func DefaultConfigPath() string {
	dir := os.Getenv("APPDATA")
	if dir != "" {
		return filepath.Join(dir, "hcloud", "cli.toml")
	}
	return ""
}
