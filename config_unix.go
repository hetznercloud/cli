// +build !windows

package cli

import (
	"os/user"
	"path/filepath"
)

func init() {
	usr, err := user.Current()
	if err != nil {
		return
	}
	if usr.HomeDir != "" {
		DefaultConfigPath = filepath.Join(usr.HomeDir, ".config", "hcloud", "cli.toml")
	}
}
