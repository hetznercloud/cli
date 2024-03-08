package util

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
)

func TestValidateExact(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return ValidateExact(&cobra.Command{
			Use: "test <arg1> <arg2>",
		}, []string{"arg1", "arg2"})
	})

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NoError(t, err)
}

func TestValidateExactMissing(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return ValidateExact(&cobra.Command{
			Use: "test <arg1> <arg2> <arg3>",
		}, []string{"arg1", "arg2"})
	})

	assert.Empty(t, stdout)
	assert.Equal(t, "test <arg1> <arg2> <arg3>\n                    ^^^^\n", stderr)
	assert.EqualError(t, err, "expected argument arg3 at position 3")
}

func TestValidateExactTooMany(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return ValidateExact(&cobra.Command{
			Use: "test <arg1> <arg2>",
		}, []string{"arg1", "arg2", "arg3"})
	})

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.EqualError(t, err, "expected exactly 2 positional arguments, but got 3")
}

func TestValidateTooMany(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return Validate(&cobra.Command{
			Use: "test <arg1> <arg2>",
		}, []string{"arg1", "arg2", "arg3"})
	})

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NoError(t, err)
}

func TestValidateExactComplexUsage(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return ValidateExact(&cobra.Command{
			Use: "test [options] [<optional-arg>] <arg1> --flag <not-an-arg> <arg2> [<arg3>]",
		}, []string{"arg1", "arg2"})
	})

	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.NoError(t, err)
}
