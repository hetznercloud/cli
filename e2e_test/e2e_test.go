//go:build e2e

package e2e

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var client = hcloud.NewClient(hcloud.WithToken(os.Getenv("HCLOUD_TOKEN")))

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
	t.Helper()
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
	t.Helper()
	cmd := newRootCommand(t)
	var buf bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOut(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

func withSuffix(s string) string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	suffix := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s", s, suffix)
}
