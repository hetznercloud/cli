//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/assertjson"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestPlacementGroup(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "placement-group", "create")
	assert.Empty(t, out)
	require.EqualError(t, err, `required flag(s) "name", "type" not set`)

	pgName := withSuffix("test-placement-group")
	pgID, err := createPlacementGroup(t, pgName, "--type", "spread")
	require.NoError(t, err)

	out, err = runCommand(t, "placement-group", "add-label", "non-existing-placement-group", "foo=bar")
	require.EqualError(t, err, "placement group not found: non-existing-placement-group")
	assert.Empty(t, out)

	out, err = runCommand(t, "placement-group", "add-label", pgName, "foo=bar")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) foo added to placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "add-label", pgName, "baz=qux")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz added to placement group %d\n", pgID), out)

	oldPgName := pgName
	pgName = withSuffix("new-test-placement-group")

	out, err = runCommand(t, "placement-group", "update", oldPgName, "--name", pgName)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("placement group %s updated\n", oldPgName), out)

	out, err = runCommand(t, "placement-group", "list", "-o=columns=id,name,servers,type,created,age")
	require.NoError(t, err)
	assert.Regexp(t, `ID +NAME +SERVERS +TYPE +CREATED +AGE
[0-9]+ +new-test-placement-group-[0-9a-f]{8} +0 servers +spread .*? (?:just now|[0-9]+s)
`, out)

	out, err = runCommand(t, "placement-group", "describe", strconv.Itoa(pgID))
	require.NoError(t, err)
	assert.Regexp(t, `^ID:\s+[0-9]+
Name:\s+new-test-placement-group-[0-9a-f]{8}
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
    "name": "%s",
    "labels": {
      "baz": "qux",
      "foo": "bar"
    },
    "created": "<ignore-diff>",
    "servers": [],
    "type": "spread"
  }
]`, pgID, pgName)), []byte(out))

	out, err = runCommand(t, "placement-group", "remove-label", pgName, "baz")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Label(s) baz removed from placement group %d\n", pgID), out)

	out, err = runCommand(t, "placement-group", "delete", strconv.Itoa(pgID))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("placement group %d deleted\n", pgID), out)
}

func createPlacementGroup(t *testing.T, name string, args ...string) (int, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.PlacementGroup.Delete(context.Background(), &hcloud.PlacementGroup{Name: name})
	})

	out, err := runCommand(t, append([]string{"placement-group", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^Placement group [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.Atoi(out[16 : len(out)-9])
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.PlacementGroup.Delete(context.Background(), &hcloud.PlacementGroup{ID: int64(id)})
	})
	return id, nil
}
