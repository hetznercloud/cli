package image_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := image.ListCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()
	fx.Client.ImageClient.EXPECT().
		AllWithOpts(
			gomock.Any(),
			hcloud.ImageListOpts{
				ListOpts:          hcloud.ListOpts{PerPage: 50},
				Sort:              []string{"id:asc"},
				IncludeDeprecated: true,
			},
		).
		Return([]*hcloud.Image{
			{
				ID:           123,
				Type:         hcloud.ImageTypeSystem,
				Name:         "test",
				Architecture: hcloud.ArchitectureX86,
				ImageSize:    20.0,
				DiskSize:     15,
				Created:      time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:           456,
				Type:         hcloud.ImageTypeSystem,
				Name:         "deprecated",
				Architecture: hcloud.ArchitectureX86,
				ImageSize:    20.0,
				DiskSize:     15,
				Created:      time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
				DeprecatableResource: hcloud.DeprecatableResource{Deprecation: &hcloud.DeprecationInfo{
					Announced:        time.Date(2036, 1, 1, 0, 0, 0, 0, time.UTC),
					UnavailableAfter: time.Date(2036, 12, 31, 0, 0, 0, 0, time.UTC),
				}},
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `ID    TYPE     NAME         DESCRIPTION   ARCHITECTURE   IMAGE SIZE   DISK SIZE   CREATED                   DEPRECATED
123   system   test         -             x86            20.00 GB     15 GB       2036-08-20 12:00:00 UTC   -
456   system   deprecated   -             x86            20.00 GB     15 GB       2036-08-20 12:00:00 UTC   2036-01-01 00:00:00 UTC
`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
