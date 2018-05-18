// +build windows

package cli

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
