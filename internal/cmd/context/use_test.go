package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestUse(t *testing.T) {

	testConfig := `active_context = "my-context"

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
		name   string
		args   []string
		config string
		err    string
		expErr string
		expOut string
	}

	testCases := []testCase{
		{
			name:   "use different context",
			args:   []string{"my-other-context"},
			config: testConfig,
			expOut: `active_context = "my-other-context"

[[contexts]]
  name = "my-context"
  token = "super secret token"

[[contexts]]
  name = "my-other-context"
  token = "another super secret token"

[[contexts]]
  name = "my-third-context"
  token = "yet another super secret token"
`,
		},
		{
			name:   "use active context",
			args:   []string{"my-context"},
			config: testConfig,
			expOut: testConfig,
		},
		{
			name:   "use non-existing context",
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

			cmd := context.NewUseCommand(fx.State())
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
