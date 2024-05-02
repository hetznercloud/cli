package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestUnset(t *testing.T) {
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
			name: "unset in current context",
			args: []string{"quiet"},
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

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name: "unset in other context",
			preRun: func() {
				// usually you would do this with a flag, but it is only defined on the root command,
				// so we can't use it here
				_ = os.Setenv("HCLOUD_CONTEXT", "other_context")
			},
			postRun: func() {
				_ = os.Unsetenv("HCLOUD_CONTEXT")
			},
			args: []string{"poll-interval"},
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
`,
		},
		{
			name: "unset globally",
			args: []string{"debug", "--global"},
			expOut: `active_context = "test_context"

[preferences]
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

			cmd := configCmd.NewUnsetCommand(fx.State())

			out, errOut, err := fx.Run(cmd, tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
