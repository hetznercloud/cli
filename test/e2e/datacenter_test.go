//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatacenter(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "datacenter", "list")
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "NAME", "DESCRIPTION", "LOCATION").Newline().
					AnyTimes(NewRegex().Int().Whitespace().Identifier().Whitespace().AnyString().Whitespace().LocationName().Newline()).
					End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			out, err := runCommand(t, "datacenter", "list", "-o=json")
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal([]byte(out), new([]any)))
		})
	})

	t.Run("describe", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "datacenter", "describe", "123456")
			require.EqualError(t, err, "datacenter not found: 123456")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "datacenter", "describe", TestDatacenterID)
			require.NoError(t, err)

			assert.Regexp(t,
				NewRegex().Start().
					Lit("ID:").Whitespace().Int().Newline().
					Lit("Name:").Whitespace().Identifier().Newline().
					Lit("Description:").Whitespace().AnyString().Newline().
					Lit("Location:").Newline().
					Lit("  Name:").Whitespace().LocationName().Newline().
					Lit("  Description:").Whitespace().AnyString().Newline().
					Lit("  Country:").Whitespace().CountryCode().Newline().
					Lit("  City:").Whitespace().AnyString().Newline().
					Lit("  Latitude:").Whitespace().Float().Newline().
					Lit("  Longitude:").Whitespace().Float().Newline().
					Lit("Server Types:").Newline().
					AnyTimes(
						NewRegex().
							Lit("  - ID: ").Int().Whitespace().
							Lit("Name: ").Identifier().Whitespace().
							Lit("Supported: ").OneOfLit("true", "false").Whitespace().
							Lit("Available: ").OneOfLit("true", "false").Newline(),
					).End(),
				out,
			)
		})
	})
}
