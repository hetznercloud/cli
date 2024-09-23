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
			assert.Regexp(t, `ID:\s+[0-9]+
Name:\s+[a-z0-9\-]+
Description:\s+[a-zA-Z0-9\- ]+
Location:
 +Name:\s+[a-z0-9]+
 +Description:\s+[a-zA-Z0-9\- ]+
 +Country:\s+[A-Z]{2}
 +City:\s+[A-Za-z]+
 +Latitude:\s+[0-9\.]+
 +Longitude:\s+[0-9\.]+
Server Types:
(\s+- ID: [0-9]+\s+Name: [a-z0-9]+\s+Supported: (true|false)\s+Available: (true|false))+
`, out)
		})
	})
}
