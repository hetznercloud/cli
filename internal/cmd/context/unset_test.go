package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestUnset(t *testing.T) {
	// Make sure we don't have any environment variables set that could interfere with the test
	t.Setenv("HCLOUD_TOKEN", "")
	t.Setenv("HCLOUD_CONTEXT", "")

	testConfig := `active_context = "my-context"

[[contexts]]
  name = "my-context"
  token = "super secret token"

[[contexts]]
  name = "my-other-context"
  token = "another super secret token"

[[contexts]]
  name = "my-third-context"
  token = "yet another super secret token"
`

	expOut := `active_context = ""

[[contexts]]
  name = "my-context"
  token = "super secret token"

[[contexts]]
  name = "my-other-context"
  token = "another super secret token"

[[contexts]]
  name = "my-third-context"
  token = "yet another super secret token"
`

	fx := testutil.NewFixtureWithConfigFile(t, []byte(testConfig))
	defer fx.Finish()

	cmd := context.NewUnsetCommand(fx.State())
	out, errOut, err := fx.Run(cmd, []string{})

	require.NoError(t, err)
	assert.Equal(t, "", errOut)
	assert.Equal(t, expOut, out)
}
