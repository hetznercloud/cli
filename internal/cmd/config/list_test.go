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
			expOut: `KEY                VALUE                       ORIGIN
context            test_context                config file
debug              yes                         environment
default-ssh-keys   [1 2 3]                     config file
endpoint           https://test-endpoint.com   flag
poll-interval      1.234s                      environment
quiet              yes                         flag
token              [redacted]                  config file
`,
		},
		{
			name: "no origin",
			args: []string{"-o=columns=key,value"},
			expOut: `KEY                VALUE
context            test_context
debug              yes
default-ssh-keys   [1 2 3]
endpoint           https://test-endpoint.com
poll-interval      1.234s
quiet              yes
token              [redacted]
`,
		},
		{
			name: "no header",
			args: []string{"-o=noheader"},
			expOut: `context            test_context                config file
debug              yes                         environment
default-ssh-keys   [1 2 3]                     config file
endpoint           https://test-endpoint.com   flag
poll-interval      1.234s                      environment
quiet              yes                         flag
token              [redacted]                  config file
`,
		},
		{
			name: "allow sensitive",
			args: []string{"--allow-sensitive"},
			expOut: `KEY                VALUE                       ORIGIN
context            test_context                config file
debug              yes                         environment
default-ssh-keys   [1 2 3]                     config file
endpoint           https://test-endpoint.com   flag
poll-interval      1.234s                      environment
quiet              yes                         flag
token              super secret token          config file
`,
		},
		{
			name: "json",
			args: []string{"-o=json"},
			expOut: `{
  "options": [
    {
      "key": "context",
      "value": "test_context",
      "origin": "config file"
    },
    {
      "key": "debug",
      "value": true,
      "origin": "environment"
    },
    {
      "key": "default-ssh-keys",
      "value": [
        "1",
        "2",
        "3"
      ],
      "origin": "config file"
    },
    {
      "key": "endpoint",
      "value": "https://test-endpoint.com",
      "origin": "flag"
    },
    {
      "key": "poll-interval",
      "value": 1234000000,
      "origin": "environment"
    },
    {
      "key": "quiet",
      "value": true,
      "origin": "flag"
    },
    {
      "key": "token",
      "value": "[redacted]",
      "origin": "config file"
    }
  ]
}
`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, "testdata/cli.toml")
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
