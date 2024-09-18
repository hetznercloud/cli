//go:build e2e

package e2e

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/assertjson"
)

func TestPlacementGroup(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "placement-group", "create")
	assert.Empty(t, out)
	require.EqualError(t, err, `required flag(s) "name", "type" not set`)

	out, err = runCommand(t, "placement-group", "create", "--name", "test-placement-group", "--type", "spread")
	require.NoError(t, err)
	if !assert.Regexp(t, `^Placement group [0-9]+ created\n$`, out) {
		return
	}

	pgID, err := strconv.Atoi(out[16 : len(out)-9])
	require.NoError(t, err)

	out, err = runCommand(t, "placement-group", "add-label", "non-existing-placement-group", "foo=bar")
	require.EqualError(t, err, "placement group not found: non-existing-placement-group")
	assert.Empty(t, out)

	out, err = runCommand(t, "placement-group", "add-label", "test-placement-group", "foo=bar")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "add-label", "test-placement-group", "baz=qux")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "update", "test-placement-group", "--name", "new-test-placement-group")
	require.NoError(t, err)
	assert.Equal(t, "placement group test-placement-group updated\n", out)

	out, err = runCommand(t, "placement-group", "list", "-o=columns=id,name,servers,type,created,age")
	require.NoError(t, err)
	assert.Regexp(t, `ID +NAME +SERVERS +TYPE +CREATED +AGE
[0-9]+ +new-test-placement-group +0 servers +spread .*? (?:just now|[0-9]+s)
`, out)

	out, err = runCommand(t, "placement-group", "describe", strconv.Itoa(pgID))
	require.NoError(t, err)
	assert.Regexp(t, `^ID:\s+[0-9]+
Name:\s+new-test-placement-group
Created:\s+.*?
Labels:
\s+(baz: qux|foo: bar)
\s+(baz: qux|foo: bar)
Servers:
Type:\s+spread
$`, out)

	out, err = runCommand(t, "placement-group", "list", "-o=json")
	require.NoError(t, err)
	assertjson.Equal(t, []byte(fmt.Sprintf(`
[
  {
    "id": %d,
    "name": "new-test-placement-group",
    "labels": {
      "baz": "qux",
      "foo": "bar"
    },
    "created": "<ignore-diff>",
    "servers": [],
    "type": "spread"
  }
]`, pgID)), []byte(out))

	out, err = runCommand(t, "placement-group", "remove-label", "new-test-placement-group", "baz")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "delete", strconv.Itoa(pgID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("placement group %d deleted\n", pgID), out)
}
