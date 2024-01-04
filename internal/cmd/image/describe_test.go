package image

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	img := &hcloud.Image{
		ID:           123,
		Type:         hcloud.ImageTypeSystem,
		Status:       hcloud.ImageStatusAvailable,
		Name:         "test",
		Created:      time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
		Description:  "Test image",
		ImageSize:    20.0,
		DiskSize:     20.0,
		Architecture: hcloud.ArchitectureX86,
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "test", hcloud.ArchitectureX86).
		Return(img, nil, nil)

	out, _, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Type:		system
Status:		available
Name:		test
Created:	%s (%s)
Description:	Test image
Image size:	20.00 GB
Disk size:	20 GB
OS flavor:	
OS version:	-
Architecture:	x86
Rapid deploy:	no
Protection:
  Delete:	no
Labels:
  key: value
`, util.Datetime(img.Created), humanize.Time(img.Created))

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
