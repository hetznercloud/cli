package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestUnset(t *testing.T) {
	os.Clearenv()

	_, deleteNestedOption := config.NewTestOption(
		"nested.option",
		"nested option",
		"foo",
		config.OptionFlagPreference,
		nil,
	)
	defer deleteNestedOption()

	_, deleteDeeplyNestedOption := config.NewTestOption(
		"deeply.nested.option",
		"deeply nested option",
		"foo",
		config.OptionFlagPreference,
		nil,
	)
	defer deleteDeeplyNestedOption()

	testConfig := `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.deeply]
      [contexts.preferences.deeply.nested]
        option = "bar"
    [contexts.preferences.nested]
      option = "foo"

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`

	type testCase struct {
		name    string
		args    []string
		config  string
		expOut  string
		expErr  string
		err     string
		preRun  func()
		postRun func()
	}

	testCases := []testCase{
		{
			name:   "unset in current context",
			args:   []string{"quiet"},
			config: testConfig,
			expOut: `Unset 'quiet' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    [contexts.preferences.deeply]
      [contexts.preferences.deeply.nested]
        option = "bar"
    [contexts.preferences.nested]
      option = "foo"

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
			args:   []string{"poll-interval"},
			config: testConfig,
			expOut: `Unset 'poll-interval' in context 'other_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.deeply]
      [contexts.preferences.deeply.nested]
        option = "bar"
    [contexts.preferences.nested]
      option = "foo"

[[contexts]]
  name = "other_context"
  token = "another super secret token"
`,
		},
		{
			name:   "unset globally",
			args:   []string{"debug", "--global"},
			config: testConfig,
			expOut: `Unset 'debug' globally
active_context = "test_context"

[preferences]
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.deeply]
      [contexts.preferences.deeply.nested]
        option = "bar"
    [contexts.preferences.nested]
      option = "foo"

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "unset non existing",
			args:   []string{"non-existing"},
			config: testConfig,
			err:    "unknown preference: non-existing",
			expErr: "Error: unknown preference: non-existing\n",
		},
		{
			name:   "unset not set",
			args:   []string{"debug-file"},
			config: testConfig,
			expErr: "Warning: key 'debug-file' was not set\n",
			expOut: `Unset 'debug-file' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.deeply]
      [contexts.preferences.deeply.nested]
        option = "bar"
    [contexts.preferences.nested]
      option = "foo"

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "unset nested",
			args:   []string{"nested.option"},
			config: testConfig,
			expOut: `Unset 'nested.option' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.deeply]
      [contexts.preferences.deeply.nested]
        option = "bar"

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "unset deeply nested",
			args:   []string{"deeply.nested.option"},
			config: testConfig,
			expOut: `Unset 'deeply.nested.option' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      option = "foo"

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

			fx := testutil.NewFixtureWithConfigFile(t, []byte(tt.config))
			defer fx.Finish()

			cmd := configCmd.NewUnsetCommand(fx.State())

			out, errOut, err := fx.Run(cmd, tt.args)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.err)
			}
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
