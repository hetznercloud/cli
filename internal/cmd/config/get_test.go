package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configCmd "github.com/hetznercloud/cli/internal/cmd/config"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestGet(t *testing.T) {
	type testCase struct {
		key    string
		args   []string
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
			key:    "default-ssh-keys",
			expOut: "[1 2 3]\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.key, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, "testdata/cli.toml")
			defer fx.Finish()

			cmd := configCmd.NewGetCommand(fx.State())

			// sets flags and env variables
			setTestValues(fx.Config)
			out, errOut, err := fx.Run(cmd, append(tt.args, tt.key))

			assert.NoError(t, err)
			assert.Equal(t, tt.expErr, errOut)
			assert.Equal(t, tt.expOut, out)
		})
	}
}
