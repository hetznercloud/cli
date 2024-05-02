package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestAdd(t *testing.T) {
	type testCase struct {
		name    string
		args    []string
		expOut  string
		expErr  string
		preRun  func()
		postRun func()
	}

	testCases := []testCase{
		{
			name: "add to existing",
			args: []string{"default-ssh-keys", "a", "b", "c"},
			expOut: `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    default_ssh_keys = ["1", "2", "3", "a", "b", "c"]
    endpoint = "https://test-endpoint.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name: "global add to empty",
			args: []string{"--global", "default-ssh-keys", "a", "b", "c"},
			expOut: `active_context = "test_context"

[preferences]
  debug = true
  default_ssh_keys = ["a", "b", "c"]
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    default_ssh_keys = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name: "global add to empty duplicate",
			args: []string{"--global", "default-ssh-keys", "c", "b", "c", "a", "a"},
			expOut: `active_context = "test_context"

[preferences]
  debug = true
  default_ssh_keys = ["a", "b", "c"]
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    default_ssh_keys = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			preRun: func() {
				_ = os.Setenv("HCLOUD_CONTEXT", "other_context")
			},
			postRun: func() {
				_ = os.Unsetenv("HCLOUD_CONTEXT")
			},
			name: "add to other context",
			args: []string{"default-ssh-keys", "I", "II", "III"},
			expOut: `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    default_ssh_keys = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    default_ssh_keys = ["I", "II", "III"]
    poll_interval = "1.234s"
`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preRun != nil {
				tt.preRun()
			}
			if tt.postRun != nil {
				defer tt.postRun()
			}

			fx := testutil.NewFixtureWithConfigFile(t, "testdata/cli.toml")
			defer fx.Finish()

			cmd := configCmd.NewAddCommand(fx.State())

			out, errOut, err := fx.Run(cmd, tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
