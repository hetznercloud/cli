package util

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		use            string
		args           []string
		exact          bool
		expectedStdout string
		expectedStderr string
		expectedErr    string
	}{
		{
			name: "correct usage",
			use:  "test <arg1> <arg2>",
			args: []string{"arg1", "arg2"},
		},
		{
			name:  "correct usage exact",
			use:   "test <arg1> <arg2>",
			args:  []string{"arg1", "arg2"},
			exact: true,
		},
		{
			name:           "missing arg",
			use:            "test <arg1> <arg2>",
			args:           []string{"arg1"},
			expectedStderr: "test <arg1> <arg2>\n             ^^^^\n",
			expectedErr:    "expected argument arg2 at position 2",
		},
		{
			name:           "missing arg exact",
			use:            "test <arg1> <arg2>",
			args:           []string{"arg1"},
			expectedStderr: "test <arg1> <arg2>\n             ^^^^\n",
			expectedErr:    "expected argument arg2 at position 2",
			exact:          true,
		},
		{
			name: "too many args",
			use:  "test <arg1> <arg2>",
			args: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:        "too many args exact",
			use:         "test <arg1> <arg2>",
			args:        []string{"arg1", "arg2", "arg3"},
			exact:       true,
			expectedErr: "expected exactly 2 positional arguments, but got 3",
		},
		{
			name: "correct usage complex",
			use:  "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args: []string{"arg1", "arg2"},
		},
		{
			name:  "correct usage complex exact",
			use:   "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:  []string{"arg1", "arg2"},
			exact: true,
		},
		{
			name:           "complex missing arg",
			use:            "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:           []string{"arg1"},
			expectedStderr: "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]\n                                                            ^^^^\n",
			expectedErr:    "expected argument arg2 at position 2",
		},
		{
			name:           "complex missing arg exact",
			use:            "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:           []string{"arg1"},
			exact:          true,
			expectedStderr: "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]\n                                                            ^^^^\n",
			expectedErr:    "expected argument arg2 at position 2",
		},
		{
			name: "complex too many args",
			use:  "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args: []string{"arg1", "arg2", "arg3"},
		},

		{
			name:        "complex too many args exact",
			use:         "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:        []string{"arg1", "arg2", "arg3"},
			exact:       true,
			expectedErr: "expected exactly 2 positional arguments, but got 3",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			stdout, stderr, err := testutil.CaptureOutStreams(func() error {
				cmd := &cobra.Command{Use: test.use}
				if test.exact {
					return ValidateExact(cmd, test.args)
				} else {
					return Validate(cmd, test.args)
				}
			})

			assert.Equal(t, test.expectedStdout, stdout)
			assert.Equal(t, test.expectedStderr, stderr)
			if test.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedErr)
			}
		})
	}
}
