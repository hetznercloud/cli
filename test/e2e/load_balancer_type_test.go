//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancerType(t *testing.T) {
	t.Parallel()

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer-type", "list")
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "NAME", "DESCRIPTION", "MAX SERVICES", "MAX CONNECTIONS", "MAX TARGETS").Newline().
					AnyTimes(NewRegex().
						Int().Whitespace().
						Identifier().Whitespace().
						AnyString().Whitespace().
						Int().Whitespace().
						Int().Whitespace().
						Int().Newline()).
					End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer-type", "list", "-o=json")
			require.NoError(t, err)
			assert.True(t, json.Valid([]byte(out)), "is valid JSON")
		})
	})

	t.Run("describe", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer-type", "describe", "non-existing")
			require.EqualError(t, err, "Load Balancer Type not found: non-existing")
			assert.Empty(t, out)
		})

		t.Run("normal", func(t *testing.T) {
			out, err := runCommand(t, "load-balancer-type", "describe", TestLoadBalancerTypeName)
			require.NoError(t, err)
			assert.Regexp(t,
				NewRegex().Start().
					Lit("ID:").Whitespace().Int().Newline().
					Lit("Name:").Whitespace().Identifier().Newline().
					Lit("Description:").Whitespace().AnyString().Newline().
					Lit("Max Services:").Whitespace().Int().Newline().
					Lit("Max Connections:").Whitespace().Int().Newline().
					Lit("Max Targets:").Whitespace().Int().Newline().
					Lit("Max assigned Certificates:").Whitespace().Int().Newline().
					Lit("Pricings per Location:").Newline().
					AnyTimes(
						NewRegex().
							Lit("  - Location:").Whitespace().LocationName().Newline().
							Lit("    Hourly:").Whitespace().Price().Newline().
							Lit("    Monthly:").Whitespace().Price().Newline().
							Lit("    Included Traffic:").Whitespace().IBytes().Newline().
							Lit("    Additional Traffic:").Whitespace().Price().Lit(" per TB").Newline().
							Newline(),
					).End(),
				out,
			)
		})
	})
}
