package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestDelete(t *testing.T) {
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
			name:   "delete active context",
			args:   []string{"my-context"},
			config: testConfig,
			expErr: "Warning: You are deleting the currently active context. Please select a new active context.\n",
			expOut: `active_context = ""

[[contexts]]
  name = "my-other-context"
  token = "super secret token"
`,
		},
		{
			name:   "delete inactive context",
			args:   []string{"my-other-context"},
			config: testConfig,
			expOut: `active_context = "my-context"

[[contexts]]
  name = "my-context"
  token = "super secret token"
`,
		},
		{
			name:   "delete non-existing context",
			args:   []string{"non-existing-context"},
			config: testConfig,
			err:    "context not found: non-existing-context",
			expErr: "Error: context not found: non-existing-context\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, []byte(tt.config))
			defer fx.Finish()

			cmd := context.NewDeleteCommand(fx.State())
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
