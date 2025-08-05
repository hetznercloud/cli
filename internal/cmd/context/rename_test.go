package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestRename(t *testing.T) {
	testConfig := `
active_context = "my-context"

[[contexts]]
name = "my-context"
token = "super secret token"

[[contexts]]
name = "my-other-context"
token = "super secret token"
`

	type testCase struct {
		name   string
		args   []string
		config string
		err    string
		expErr string
		expOut string
	}

	testCases := []testCase{
		{
			name:   "rename active context",
			args:   []string{"my-context", "my-renamed-context"},
			config: testConfig,
			expOut: `active_context = "my-renamed-context"

[[contexts]]
  name = "my-renamed-context"
  token = "super secret token"

[[contexts]]
  name = "my-other-context"
  token = "super secret token"
`,
		},
		{
			name:   "rename inactive context",
			args:   []string{"my-other-context", "my-other-renamed-context"},
			config: testConfig,
			expOut: `active_context = "my-context"

[[contexts]]
  name = "my-context"
  token = "super secret token"

[[contexts]]
  name = "my-other-renamed-context"
  token = "super secret token"
`,
		},
		{
			name:   "rename non-existing context",
			args:   []string{"non-existing-context", "non-existing-renamed-context"},
			config: testConfig,
			err:    "context not found: non-existing-context",
			expErr: "Error: context not found: non-existing-context\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, []byte(tt.config))
			defer fx.Finish()

			cmd := context.NewRenameCommand(fx.State())
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
