package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestAdd(t *testing.T) {
	os.Clearenv()

	_, deleteArrayOption := config.NewTestOption[[]string](
		"array-option",
		"array option",
		nil,
		config.OptionFlagPreference,
		nil,
	)
	defer deleteArrayOption()

	_, deleteNestedArrayOption := config.NewTestOption[[]string](
		"nested.array-option",
		"nested array option",
		nil,
		config.OptionFlagPreference,
		nil,
	)
	defer deleteNestedArrayOption()

	testConfig := `active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      array_option = ["1", "2", "3"]

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
		preRun  func()
		postRun func()
	}

	testCases := []testCase{
		{
			name:   "add to existing",
			args:   []string{"array-option", "a", "b", "c"},
			config: testConfig,
			expOut: `Added '[a b c]' to 'array-option' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["1", "2", "3", "a", "b", "c"]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      array_option = ["1", "2", "3"]

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "add to nested",
			args:   []string{"nested.array-option", "a", "b", "c"},
			config: testConfig,
			expOut: `Added '[a b c]' to 'nested.array-option' in context 'test_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      array_option = ["1", "2", "3", "a", "b", "c"]

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "global add to empty",
			args:   []string{"--global", "array-option", "a", "b", "c"},
			config: testConfig,
			expOut: `Added '[a b c]' to 'array-option' globally
active_context = "test_context"

[preferences]
  array_option = ["a", "b", "c"]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      array_option = ["1", "2", "3"]

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`,
		},
		{
			name:   "global add to empty duplicate",
			args:   []string{"--global", "array-option", "c", "b", "c", "a", "a"},
			config: testConfig,
			expErr: "Warning: some values were already present or duplicate\n",
			expOut: `Added '[c b a]' to 'array-option' globally
active_context = "test_context"

[preferences]
  array_option = ["c", "b", "a"]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      array_option = ["1", "2", "3"]

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
			name:   "add to other context",
			args:   []string{"array-option", "I", "II", "III"},
			config: testConfig,
			expOut: `Added '[I II III]' to 'array-option' in context 'other_context'
active_context = "test_context"

[preferences]
  debug = true
  poll_interval = "1.234s"

[[contexts]]
  name = "test_context"
  token = "super secret token"
  [contexts.preferences]
    array_option = ["1", "2", "3"]
    endpoint = "https://test-endpoint.com"
    quiet = true
    [contexts.preferences.nested]
      array_option = ["1", "2", "3"]

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    array_option = ["I", "II", "III"]
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

			cmd := configCmd.NewAddCommand(fx.State())

			out, errOut, err := fx.Run(cmd, tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
