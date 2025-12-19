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

func TestSet(t *testing.T) {
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

	_, deleteArrayOption := config.NewTestOption[[]string](
		"array-option",
		"array option",
		nil,
		config.OptionFlagPreference,
		nil,
	)
	defer deleteArrayOption()

	testConfig := `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://example.com"
    quiet = true

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
			name:   "set in current context",
			args:   []string{"debug-file", "debug.log"},
			config: testConfig,
			expOut: `Set 'debug-file' to 'debug.log' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    debug_file = "debug.log"
    endpoint = "https://example.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name: "set in other context",
			preRun: func() {
				// usually you would do this with a flag, but it is only defined on the root command,
				// so we can't use it here
				_ = os.Setenv("HCLOUD_CONTEXT", "other_context")
			},
			postRun: func() {
				_ = os.Unsetenv("HCLOUD_CONTEXT")
			},
			args:   []string{"debug", "false"},
			config: testConfig,
			expOut: `Set 'debug' to 'false' in context 'other_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://example.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    debug = false
    poll_interval = "1.234s"
`,
		},
		{
			name:   "set globally",
			args:   []string{"--global", "poll-interval", "50ms"},
			config: testConfig,
			expOut: `Set 'poll-interval' to '50ms' globally
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "50ms"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://example.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "set nested",
			args:   []string{"nested.option", "bar"},
			config: testConfig,
			expOut: `Set 'nested.option' to 'bar' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://example.com"
    quiet = true
    [contexts.preferences.nested]
      option = "bar"

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "set deeply nested",
			args:   []string{"deeply.nested.option", "bar"},
			config: testConfig,
			expOut: `Set 'deeply.nested.option' to 'bar' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    endpoint = "https://example.com"
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
			name:   "set array option",
			args:   []string{"array-option", "a", "b", "c"},
			config: testConfig,
			expOut: `Set 'array-option' to '[a b c]' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["a", "b", "c"]
    endpoint = "https://example.com"
    quiet = true

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "set non existing",
			args:   []string{"non-existing", "value"},
			config: testConfig,
			expErr: "Error: unknown preference: non-existing\n",
			err:    "unknown preference: non-existing",
		},
		{
			name:   "set non-preference",
			args:   []string{"token", "value"},
			config: testConfig,
			expErr: "Error: unknown preference: token\n",
			err:    "unknown preference: token",
		},
		{
			name: "set in empty config global",
			args: []string{"debug", "false", "--global"},
			expOut: `Set 'debug' to 'false' globally
active_context = ""

[preferences]
  debug = false
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

			cmd := configCmd.NewSetCommand(fx.State())

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
