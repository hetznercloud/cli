package image

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "test", hcloud.ArchitectureX86).
		Return(&hcloud.Image{
			ID:           123,
			Type:         hcloud.ImageTypeSystem,
			Status:       hcloud.ImageStatusAvailable,
			Name:         "test",
			Created:      time.Date(1905, 10, 6, 12, 0, 0, 0, time.UTC),
			Description:  "Test image",
			ImageSize:    20.0,
			DiskSize:     20.0,
			Architecture: hcloud.ArchitectureX86,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Type:		system
Status:		available
Name:		test
Created:	Fri Oct  6 12:00:00 UTC 1905 (a long while ago)
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
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
