package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestList(t *testing.T) {
	os.Clearenv()

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
		name   string
		args   []string
		expOut string
		expErr string
	}

	testCases := []testCase{
		{
			name: "default",
			args: []string{},
			expOut: `KEY                    VALUE
context                test_context
debug                  yes
deeply.nested.option   bar
endpoint               https://test-endpoint.com
poll-interval          1.234s
quiet                  yes
token                  [redacted]
`,
		},
		{
			name: "only key",
			args: []string{"-o=columns=key"},
			expOut: `KEY
context
debug
deeply.nested.option
endpoint
poll-interval
quiet
token
`,
		},
		{
			name: "no header",
			args: []string{"-o=noheader"},
			expOut: `context                test_context
debug                  yes
deeply.nested.option   bar
endpoint               https://test-endpoint.com
poll-interval          1.234s
quiet                  yes
token                  [redacted]
`,
		},
		{
			name: "allow sensitive",
			args: []string{"--allow-sensitive"},
			expOut: `KEY                    VALUE
context                test_context
debug                  yes
deeply.nested.option   bar
endpoint               https://test-endpoint.com
poll-interval          1.234s
quiet                  yes
token                  super secret token
`,
		},
		{
			name: "json",
			args: []string{"-o=json"},
			expOut: `{
  "options": [
    {
      "key": "context",
      "value": "test_context"
    },
    {
      "key": "debug",
      "value": true
    },
    {
      "key": "deeply.nested.option",
      "value": "bar"
    },
    {
      "key": "endpoint",
      "value": "https://test-endpoint.com"
    },
    {
      "key": "poll-interval",
      "value": 1234000000
    },
    {
      "key": "quiet",
      "value": true
    },
    {
      "key": "token",
      "value": "[redacted]"
    }
  ]
}
`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, []byte(testConfig))
			defer fx.Finish()

			cmd := configCmd.NewListCommand(fx.State())

			setTestValues(fx.Config)
			out, errOut, err := fx.Run(cmd, tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}

func setTestValues(cfg config.Config) {
	_ = os.Setenv("HCLOUD_POLL_INTERVAL", "1234ms")
	_ = os.Setenv("HCLOUD_DEBUG", "true")
	_ = cfg.FlagSet().Set("endpoint", "https://test-endpoint.com")
	_ = cfg.FlagSet().Set("quiet", "true")
}
