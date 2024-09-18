//go:build e2e

package e2e_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func TestMain(m *testing.M) {
	tok := os.Getenv("HCLOUD_TOKEN")
	if tok == "" {
		fmt.Println("HCLOUD_TOKEN is not set")
		os.Exit(1)
		return
	}
	os.Exit(m.Run())
}

func newRootCommand(t *testing.T) *cobra.Command {
	cfg := config.New()
	if err := cfg.Read("config.toml"); err != nil {
		t.Fatalf("unable to read config file \"%s\": %s\n", cfg.Path(), err)
	}

	s, err := state.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	return cli.NewRootCommand(s)
}

func runCommand(t *testing.T, args ...string) (string, error) {
	cmd := newRootCommand(t)
	var buf bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOut(&buf)
	err := cmd.Execute()
	return buf.String(), err
}
