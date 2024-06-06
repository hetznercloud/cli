package context_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestActive(t *testing.T) {

	testConfig := `
active_context = "my-context"

[[contexts]]
name = "my-context"
token = "super secret token"
`

	type testCase struct {
		name    string
		args    []string
		config  string
		err     string
		expOut  string
		expErr  string
		preRun  func()
		postRun func()
	}

	testCases := []testCase{
		{
			name:   "no arguments",
			args:   []string{},
			config: testConfig,
			expOut: "my-context\n",
		},
		{
			name: "no config",
			args: []string{},
		},
		{
			name:   "from env",
			args:   []string{},
			config: testConfig,
			preRun: func() {
				_ = os.Setenv("HCLOUD_CONTEXT", "abcdef")
			},
			postRun: func() {
				_ = os.Unsetenv("HCLOUD_CONTEXT")
			},
			// 'abcdef' does not exist, so there is nothing printed to stdout.
			// The warning 'active context "abcdef" not found' should be printed to stderr during config loading, which
			// is before stderr is captured.
		},
		{
			name:   "invalid config",
			args:   []string{},
			config: `active_context = "invalid-context-name"`,
			// if there is no context with the name of the active_context, there should be no output. See above
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

			cmd := context.NewActiveCommand(fx.State())
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
