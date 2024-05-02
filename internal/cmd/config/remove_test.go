package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestRemove(t *testing.T) {
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
			name: "remove from existing",
			args: []string{"default-ssh-keys", "2", "3"},
			expOut: `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    default_ssh_keys = ["1"]
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
			name: "remove all from existing",
			args: []string{"default-ssh-keys", "1", "2", "3"},
			expOut: `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
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

			cmd := configCmd.NewRemoveCommand(fx.State())

			out, errOut, err := fx.Run(cmd, tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
