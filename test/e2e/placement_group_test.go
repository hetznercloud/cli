//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	t.Run("add-label", func(t *testing.T) {
		t.Run("non-existing-placement-group", func(t *testing.T) {
			out, err = runCommand(t, "placement-group", "add-label", "non-existing-placement-group", "foo=bar")
			require.EqualError(t, err, "Placement Group not found: non-existing-placement-group")
			assert.Empty(t, out)
		})

		t.Run("1", func(t *testing.T) {
			out, err = runCommand(t, "placement-group", "add-label", pgName, "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to Placement Group %d\n", pgID), out)
		})

		t.Run("2", func(t *testing.T) {
			out, err = runCommand(t, "placement-group", "add-label", pgName, "baz=qux")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz added to Placement Group %d\n", pgID), out)
		})
	})

	oldPgName := pgName
	pgName = withSuffix("new-test-placement-group")

	t.Run("update-name", func(t *testing.T) {
		out, err := runCommand(t, "placement-group", "update", oldPgName, "--name", pgName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Placement Group %s updated\n", oldPgName), out)
	})

	t.Run("list", func(t *testing.T) {
		t.Run("table", func(t *testing.T) {
			out, err := runCommand(t, "placement-group", "list", "-o=columns=id,name,servers,type,created,age")
			require.NoError(t, err)

			assert.Regexp(t,
				NewRegex().Start().
					SeparatedByWhitespace("ID", "NAME", "SERVERS", "TYPE", "CREATED", "AGE").OptionalWhitespace().Newline().
					Int().Whitespace().
					Raw(`new-test-placement-group-[0-9a-f]{8}`).Whitespace().
					Lit("0 servers").Whitespace().
					Lit("spread").Whitespace().
					Datetime().Whitespace().
					Age().OptionalWhitespace().Newline().End(),
				out,
			)
		})

		t.Run("json", func(t *testing.T) {
			out, err := runCommand(t, "placement-group", "list", "-o=json")
			require.NoError(t, err)
			assert.JSONEq(t, fmt.Sprintf(`
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
]`, pgID, pgName), out)
		})
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "placement-group", "describe", strconv.FormatInt(pgID, 10))
		require.NoError(t, err)

		assert.Regexp(t,
			NewRegex().Start().
				Lit("ID:").Whitespace().Int().Newline().
				Lit("Name:").Whitespace().Raw(`new-test-placement-group-[0-9a-f]{8}`).Newline().
				Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
				Lit("Type:").Whitespace().Lit("spread").Newline().
				Newline().
				Lit("Labels:").Newline().
				Lit("  baz:").Whitespace().Lit("qux").Newline().
				Lit("  foo:").Whitespace().Lit("bar").Newline().
				Newline().
				Lit("Servers:").Newline().
				Lit("  No servers").Newline().
				End(),
			out,
		)
	})

	t.Run("remove-label", func(t *testing.T) {
		out, err := runCommand(t, "placement-group", "remove-label", pgName, "baz")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Placement Group %d\n", pgID), out)
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "placement-group", "delete", strconv.FormatInt(pgID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Placement Group %d deleted\n", pgID), out)
	})
}

func createPlacementGroup(t *testing.T, name string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.PlacementGroup.Delete(context.Background(), &hcloud.PlacementGroup{Name: name})
	})

	out, err := runCommand(t, append([]string{"placement-group", "create", "--name", name}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^Placement Group [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[16:len(out)-9], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.PlacementGroup.Delete(context.Background(), &hcloud.PlacementGroup{ID: id})
	})
	return id, nil
}
