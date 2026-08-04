package cmpl_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/registration"
)

func TestRegistrationError(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	child := &cobra.Command{Use: "child"}
	root.AddCommand(child)

	cmpl.RegisterFlagCompletion(child, "missing", cmpl.SuggestNothing())

	err := registration.Error(root)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "root child")
	assert.Contains(t, err.Error(), `flag "missing"`)
}

func TestRegistrationErrorIsNilForValidRegistration(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("format", "", "format")

	cmpl.RegisterFlagCompletion(cmd, "format", cmpl.SuggestNothing())

	assert.NoError(t, registration.Error(cmd))
}
