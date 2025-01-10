//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestISO(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "iso", "list")
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "NAME", "DESCRIPTION", "TYPE", "ARCHITECTURE").Newline().
					AnyTimes(NewRegex().
						Int().Whitespace().
						Identifier().Whitespace().
						AnyString().Whitespace().
						OneOfLit("public", "private").Whitespace().
						OneOf("arm", "x86").Newline()).
					End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			out, err := runCommand(t, "iso", "list", "-o=json")
			require.NoError(t, err)
			assert.True(t, json.Valid([]byte(out)), "is valid JSON")
		})
	})

	t.Run("describe", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "iso", "describe", "non-existing-iso")
			require.EqualError(t, err, "iso not found: non-existing-iso")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "iso", "describe", TestISOName)
			require.NoError(t, err)

			assert.Regexp(t,
				NewRegex().Start().
					Lit("ID:").Whitespace().Int().Newline().
					Lit("Name:").Whitespace().Identifier().Newline().
					Lit("Description:").Whitespace().AnyString().Newline().
					Lit("Type:").Whitespace().OneOfLit("public", "private").Newline().
					Lit("Architecture:").Whitespace().OneOfLit("arm", "x86").Newline().
					End(),
				out,
			)
		})
	})
}
