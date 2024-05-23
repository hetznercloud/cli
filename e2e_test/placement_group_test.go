package e2e_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"
)

func TestPlacementGroup(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "placement-group", "create")
	assert.Empty(t, out)
	assert.EqualError(t, err, `required flag(s) "name", "type" not set`)

	out, err = runCommand(t, "placement-group", "create", "--name", "test-placement-group", "--type", "spread")
	assert.NoError(t, err)
	if !assert.Regexp(t, `^Placement group [0-9]+ created\n$`, out) {
		return
	}

	pgID, err := strconv.Atoi(out[16 : len(out)-9])
	assert.NoError(t, err)

	out, err = runCommand(t, "placement-group", "add-label", "non-existing-placement-group", "foo=bar")
	assert.EqualError(t, err, "placement group not found: non-existing-placement-group")
	assert.Empty(t, out)

	out, err = runCommand(t, "placement-group", "add-label", "test-placement-group", "foo=bar")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "add-label", "test-placement-group", "baz=qux")
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "update", "test-placement-group", "--name", "new-test-placement-group")
	assert.NoError(t, err)
	assert.Equal(t, "placement group test-placement-group updated\n", out)

	out, err = runCommand(t, "placement-group", "list", "-o=columns=id,name,servers,type,created,age")
	assert.NoError(t, err)
	assert.Regexp(t, `ID +NAME +SERVERS +TYPE +CREATED +AGE
[0-9]+ +new-test-placement-group +0 servers +spread .*? (?:just now|[0-9]+s)
`, out)

	out, err = runCommand(t, "placement-group", "describe", strconv.Itoa(pgID))
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "delete", strconv.Itoa(pgID))
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("placement group %d deleted\n", pgID), out)
}
