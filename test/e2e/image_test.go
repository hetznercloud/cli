//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "image", "list")
			require.NoError(t, err)

			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "TYPE", "NAME", "DESCRIPTION", "ARCHITECTURE", "IMAGE SIZE", "DISK SIZE", "CREATED", "DEPRECATED").OptionalWhitespace().Newline().
					AnyTimes(
						NewRegex().
							Int().Whitespace().
							OneOfLit("system", "app", "snapshot", "backup").Whitespace().
							Identifier().Whitespace().
							AnyString().Whitespace().
							OneOfLit("x86", "arm").Whitespace().
							OneOf("-", NewRegex().FileSize()).Whitespace().
							FileSize().Whitespace().
							Datetime().Whitespace().
							OneOf("-", NewRegex().Datetime()).
							OptionalWhitespace().Newline(),
					).
					End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			out, err := runCommand(t, "image", "list", "-o=json")
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal([]byte(out), new([]any)))
		})
	})

	t.Run("describe", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "image", "describe", "non-existing-image", "--architecture=x86")
			require.EqualError(t, err, "Image not found: non-existing-image")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "image", "describe", TestImageID)
			require.NoError(t, err)

			assert.Regexp(t,
				NewRegex().Start().
					Lit("ID:").Whitespace().Int().Newline().
					Lit("Type:").Whitespace().OneOfLit("system", "app", "snapshot", "backup").Newline().
					Lit("Status:").Whitespace().OneOfLit("available", "creating", "unavailable").Newline().
					Lit("Name:").Whitespace().Identifier().Newline().
					Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
					Lit("Description:").Whitespace().AnyString().Newline().
					Lit("Image size:").Whitespace().OneOf("-", NewRegex().FileSize()).Newline().
					Lit("Disk size:").Whitespace().FileSize().Newline().
					Lit("OS flavor:").Whitespace().Identifier().Newline().
					Lit("OS version:").Whitespace().Identifier().Newline().
					Lit("Architecture:").Whitespace().OneOfLit("x86", "arm").Newline().
					Lit("Rapid deploy:").Whitespace().OneOfLit("yes", "no").Newline().
					Lit("Protection:").Newline().
					Lit("  Delete:").Whitespace().OneOfLit("yes", "no").Newline().
					Lit("Labels:").Newline().
					Lit("  No labels").Newline().
					End(),
				out,
			)
		})
	})
}
