package version_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/version"
	"github.com/hetznercloud/cli/internal/testutil"
	version2 "github.com/hetznercloud/cli/internal/version"
)

func TestVersion(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := version.NewCommand(fx.State())
	out, errOut, err := fx.Run(cmd, []string{})

	expOut := fmt.Sprintf("hcloud %s\n", version2.Version)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestVersionLong(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := version.NewCommand(fx.State())
	out, errOut, err := fx.Run(cmd, []string{"--long"})

	require.NoError(t, err)
	assert.Empty(t, errOut)
	require.Regexp(t, `^hcloud .*

go version: *go1.[0-9]+.[0-9]+.*
platform: *.+/.+
revision: *(unknown|[0-9a-f]+)( \(modified\))?
revision date: *.+`, out)
}
