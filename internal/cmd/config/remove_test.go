package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestRemove(t *testing.T) {
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

	testConfig := `active_context = 'test_context'

[preferences]
debug = true
poll_interval = 1234000000

[[contexts]]
name = 'test_context'
token = 'super secret token'

[contexts.preferences]
array_option = ['1', '2', '3']
endpoint = 'https://test-endpoint.com'
quiet = true

[contexts.preferences.nested]
array_option = ['1', '2', '3']

[[contexts]]
name = 'other_context'
token = 'another super secret token'

[contexts.preferences]
poll_interval = 1234000000
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
			name:   "remove from existing",
			args:   []string{"array-option", "2", "3"},
			config: testConfig,
			expOut: `Removed '[2 3]' from 'array-option' in context 'test_context'
active_context = 'test_context'

[preferences]
debug = true
poll_interval = 1234000000

[[contexts]]
name = 'test_context'
token = 'super secret token'

[contexts.preferences]
array_option = ['1']
endpoint = 'https://test-endpoint.com'
quiet = true

[contexts.preferences.nested]
array_option = ['1', '2', '3']

[[contexts]]
name = 'other_context'
token = 'another super secret token'

[contexts.preferences]
poll_interval = 1234000000
`,
		},
		{
			name:   "remove all from existing",
			args:   []string{"array-option", "1", "2", "3"},
			config: testConfig,
			expOut: `Removed '[1 2 3]' from 'array-option' in context 'test_context'
active_context = 'test_context'

[preferences]
debug = true
poll_interval = 1234000000

[[contexts]]
name = 'test_context'
token = 'super secret token'

[contexts.preferences]
endpoint = 'https://test-endpoint.com'
quiet = true

[contexts.preferences.nested]
array_option = ['1', '2', '3']

[[contexts]]
name = 'other_context'
token = 'another super secret token'

[contexts.preferences]
poll_interval = 1234000000
`,
		},
		{
			name:   "remove from non-existing",
			args:   []string{"i-do-not-exist", "1", "2", "3"},
			config: testConfig,
			err:    "unknown preference: i-do-not-exist",
			expErr: "Error: unknown preference: i-do-not-exist\n",
		},
		{
			name:   "remove from nested",
			args:   []string{"nested.array-option", "2", "3"},
			config: testConfig,
			expOut: `Removed '[2 3]' from 'nested.array-option' in context 'test_context'
active_context = 'test_context'

[preferences]
debug = true
poll_interval = 1234000000

[[contexts]]
name = 'test_context'
token = 'super secret token'

[contexts.preferences]
array_option = ['1', '2', '3']
endpoint = 'https://test-endpoint.com'
quiet = true

[contexts.preferences.nested]
array_option = ['1']

[[contexts]]
name = 'other_context'
token = 'another super secret token'

[contexts.preferences]
poll_interval = 1234000000
`,
		},
		{
			name:   "remove all from nested",
			args:   []string{"nested.array-option", "1", "2", "3"},
			config: testConfig,
			expOut: `Removed '[1 2 3]' from 'nested.array-option' in context 'test_context'
active_context = 'test_context'

[preferences]
debug = true
poll_interval = 1234000000

[[contexts]]
name = 'test_context'
token = 'super secret token'

[contexts.preferences]
array_option = ['1', '2', '3']
endpoint = 'https://test-endpoint.com'
quiet = true

[[contexts]]
name = 'other_context'
token = 'another super secret token'

[contexts.preferences]
poll_interval = 1234000000
`,
		},
		{
			name:   "remove from non-list",
			args:   []string{"debug", "true"},
			config: testConfig,
			err:    "debug is not a list",
			expErr: "Error: debug is not a list\n",
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

			cmd := configCmd.NewRemoveCommand(fx.State())

			out, errOut, err := fx.Run(cmd, tt.args)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
