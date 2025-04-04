//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
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
						OneOf("arm", "x86").OptionalWhitespace().Newline()).
					End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			var schemas []schema.ISO
			isos, err := client.ISO.All(context.Background())
			require.NoError(t, err)
			for _, iso := range isos {
				schemas = append(schemas, hcloud.SchemaFromISO(iso))
			}
			expectedJson, err := json.Marshal(schemas)
			require.NoError(t, err)

			out, err := runCommand(t, "iso", "list", "-o=json")
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedJson), out)
		})
	})

	t.Run("describe", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "iso", "describe", "non-existing-iso")
			require.EqualError(t, err, "ISO not found: non-existing-iso")
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

		t.Run("json", func(t *testing.T) {
			iso, _, err := client.ISO.GetByName(context.Background(), TestISOName)
			require.NoError(t, err)
			expectedJson, err := json.Marshal(hcloud.SchemaFromISO(iso))
			require.NoError(t, err)

			out, err := runCommand(t, "iso", "describe", TestISOName, "-o=json")
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedJson), out)
		})
	})
}
