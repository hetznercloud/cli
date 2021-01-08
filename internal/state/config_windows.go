// +build windows

package state

import (
	"os"
	"path/filepath"
)

func init() {
	dir := os.Getenv("APPDATA")
	if dir != "" {
		DefaultConfigPath = filepath.Join(dir, "hcloud", "cli.toml")
	}
}
