//go:build !windows

package state

import (
	"os/user"
	"path/filepath"
)

func DefaultConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	if usr.HomeDir != "" {
		return filepath.Join(usr.HomeDir, ".config", "hcloud", "cli.toml")
	}
	return ""
}
