package config

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func RetrieveToken(c Config) (string, error) {
	tok, err := OptionToken.Get(c)
	if err != nil {
		return "", err
	}
	if tok != "" {
		return tok, nil
	}

	cmdStr, err := OptionTokenCommand.Get(c)
	if err != nil {
		return "", err
	}
	if cmdStr != "" {
		return tokenFromCommand(c, cmdStr)
	}
	return "", nil
}

// We might need to run the command to retrieve HCLOUD_TOKEN multiple times, for example when running tests.
// Since the command is provided by the user and might be expensive, we cache the command result after the first successful run.
var cmdCache = make(map[string]string)

func tokenFromCommand(c Config, cmdStr string) (string, error) {
	if tok, ok := cmdCache[cmdStr]; ok {
		return tok, nil
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	cmd.Env = append(cmd.Environ(), fmt.Sprintf("HCLOUD_CONTEXT=%s", c.ActiveContext().Name()))
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not retrieve token: %w", err)
	}
	tok := strings.TrimSpace(string(out))
	cmdCache[cmdStr] = tok
	return tok, nil
}
