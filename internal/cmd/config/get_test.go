package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestGet(t *testing.T) {

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

[[contexts]]
  name = "other_context"
  token = "another super secret token"
  [contexts.preferences]
    poll_interval = "1.234s"
`

	type testCase struct {
		key    string
		args   []string
		err    string
		expOut string
		expErr string
	}

	testCases := []testCase{
		{
			key:    "context",
			expOut: "test_context\n",
		},
		{
			key:    "debug",
			expOut: "true\n",
		},
		{
			key:    "endpoint",
			expOut: "https://test-endpoint.com\n",
		},
		{
			key:    "poll-interval",
			expOut: "1.234s\n",
		},
		{
			key:    "deeply.nested.option",
			expOut: "bar\n",
		},
		{
			key:    "non-existing-key",
			err:    "unknown key: non-existing-key",
			expErr: "Error: unknown key: non-existing-key\n",
		},
		{
			key:    "token",
			err:    "'token' is sensitive. use --allow-sensitive to show the value",
			expErr: "Error: 'token' is sensitive. use --allow-sensitive to show the value\n",
		},
		{
			key:    "token",
			args:   []string{"--allow-sensitive"},
			expOut: "super secret token\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.key, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, []byte(testConfig))
			defer fx.Finish()

			cmd := configCmd.NewGetCommand(fx.State())

			// sets flags and env variables
			setTestValues(fx.Config)
			out, errOut, err := fx.Run(cmd, append(tt.args, tt.key))

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
