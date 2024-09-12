package context_test

import (
	"io"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestCreate(t *testing.T) {
	// Make sure we don't have any environment variables set that could interfere with the test
	t.Setenv("HCLOUD_TOKEN", "")
	t.Setenv("HCLOUD_CONTEXT", "")

	testConfig := `
active_context = "my-context"

[[contexts]]
name = "my-context"
token = "super secret token"
`

	type testCase struct {
		name   string
		args   []string
		config string
		isTerm bool
		token  string
		err    string
		expErr string
		expOut string
	}

	testCases := []testCase{
		{
			name:   "new context",
			args:   []string{"new-context"},
			isTerm: true,
			config: testConfig,
			token:  "q4acIB6pq2CwsPqF+dNR2B6NTrv4yxmsspvDC1a02OqfMQeCz7nOk4A3pcJha8ix",
			expOut: `Token: 
active_context = "new-context"

[[contexts]]
  name = "my-context"
  token = "super secret token"

[[contexts]]
  name = "new-context"
  token = "q4acIB6pq2CwsPqF+dNR2B6NTrv4yxmsspvDC1a02OqfMQeCz7nOk4A3pcJha8ix"
Context new-context created and activated
`,
		},
		{
			name:   "not terminal",
			args:   []string{"new-context"},
			isTerm: false,
			config: testConfig,
			err:    "context create is an interactive command",
			expErr: "Error: context create is an interactive command\n",
		},
		{
			name:   "existing context",
			args:   []string{"my-context"},
			isTerm: true,
			config: testConfig,
			token:  "q4acIB6pq2CwsPqF+dNR2B6NTrv4yxmsspvDC1a02OqfMQeCz7nOk4A3pcJha8ix",
			err:    "name already used",
			expErr: "Error: name already used\n",
		},
		{
			name:   "invalid name",
			args:   []string{""},
			isTerm: true,
			config: testConfig,
			err:    "invalid name",
			expErr: "Error: invalid name\n",
		},
		{
			name:   "token too short",
			args:   []string{"new-context"},
			isTerm: true,
			config: testConfig,
			token:  "abc",
			err:    "EOF",
			expErr: "Error: EOF\n",
			expOut: "Token: \nEntered token is invalid (must be exactly 64 characters long)\nToken: \n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fx := testutil.NewFixtureWithConfigFile(t, []byte(tt.config))
			defer fx.Finish()

			fx.Terminal.EXPECT().StdoutIsTerminal().Return(tt.isTerm)

			isFirstCall := true
			fx.Terminal.EXPECT().ReadPassword(int(syscall.Stdin)).DoAndReturn(func(_ int) ([]byte, error) {
				if isFirstCall {
					isFirstCall = false
					return []byte(tt.token), nil
				}
				// return EOF after first call to prevent infinite loop
				return nil, io.EOF
			}).AnyTimes()

			cmd := context.NewCreateCommand(fx.State())
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
