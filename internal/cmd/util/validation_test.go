package util_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		use            string
		args           []string
		lenient        bool
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
			name:           "missing arg",
			use:            "test <arg1> <arg2>",
			args:           []string{"arg1"},
			expectedStderr: "test <arg1> <arg2>\n             ^^^^\n",
			expectedErr:    "expected argument(s) arg2 at position 2",
		},
		{
			name:           "missing arg variadic",
			use:            "test <arg1> <arg2>...",
			args:           []string{"arg1"},
			expectedStderr: "test <arg1> <arg2>...\n             ^^^^\n",
			expectedErr:    "expected argument(s) arg2 at position 2",
		},
		{
			name:           "too many args",
			use:            "test <arg1> <arg2>",
			args:           []string{"arg1", "arg2", "arg3"},
			expectedStderr: "test <arg1> <arg2>\n                   ^\n",
			expectedErr:    "expected exactly 2 positional argument(s), but got 3",
		},
		{
			name: "complex correct usage",
			use:  "test [options] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args: []string{"arg1", "arg2"},
		},
		{
			name: "complex correct usage variadic",
			use:  "test [options] <arg1> --flag <not-an-arg> <arg2>... [<arg3>]",
			args: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:           "complex missing arg",
			use:            "test [options] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:           []string{"arg1"},
			expectedStderr: "test [options] <arg1> --flag <not-an-arg> <arg2> [<arg3>]\n                                           ^^^^\n",
			expectedErr:    "expected argument(s) arg2 at position 2",
		},
		{
			name:           "complex missing arg variadic",
			use:            "test [options] <arg1> --flag <not-an-arg> <arg2>... [<arg3>]",
			args:           []string{"arg1"},
			expectedStderr: "test [options] <arg1> --flag <not-an-arg> <arg2>... [<arg3>]\n                                           ^^^^\n",
			expectedErr:    "expected argument(s) arg2 at position 2",
		},
		{ // note: ValidateLenient should be used here, because there are optional positional arguments; this is just for testing
			name:           "complex too many args not lenient",
			use:            "test [options] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:           []string{"arg1", "arg2", "arg3"},
			expectedStderr: "test [options] <arg1> --flag <not-an-arg> <arg2> [<arg3>]\n                                                          ^\n",
			expectedErr:    "expected exactly 2 positional argument(s), but got 3",
		},
		{
			name:    "complex too many args lenient",
			use:     "test [options] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
			args:    []string{"arg1", "arg2", "arg3"},
			lenient: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			stdout, stderr, err := testutil.CaptureOutStreams(func() error {
				cmd := &cobra.Command{Use: test.use}
				if test.lenient {
					return util.ValidateLenient(cmd, test.args)
				}
				return util.Validate(cmd, test.args)
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
