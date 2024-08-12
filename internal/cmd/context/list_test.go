package context_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestList(t *testing.T) {
	testConfig := `
active_context = "my-context"

[[contexts]]
name = "my-context"
token = "super secret token"

[[contexts]]
name = "my-other-context"
token = "another super secret token"

[[contexts]]
name = "my-third-context"
token = "yet another super secret token"
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
			name:   "default",
			args:   []string{},
			config: testConfig,
			expOut: `ACTIVE   NAME
*        my-context
         my-other-context
         my-third-context
`,
		},
		{
			name:   "no config",
			args:   []string{},
			expOut: "ACTIVE   NAME\n",
		},
		{
			name:   "no header",
			args:   []string{"-o=noheader"},
			config: testConfig,
			expOut: `*   my-context
    my-other-context
    my-third-context
`,
		},
		{
			name:   "no header only name",
			args:   []string{"-o=noheader", "-o=columns=name"},
			config: testConfig,
			expOut: `my-context
my-other-context
my-third-context
`,
		},
		{
			name:   "different context",
			args:   []string{},
			config: testConfig,
			preRun: func() {
				_ = os.Setenv("HCLOUD_CONTEXT", "my-other-context")
			},
			postRun: func() {
				_ = os.Unsetenv("HCLOUD_CONTEXT")
			},
			expOut: `ACTIVE   NAME
         my-context
*        my-other-context
         my-third-context
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

			cmd := context.NewListCommand(fx.State())
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
