package server

import (
	"context"
	"testing"
	"time"

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

	cmd := DescribeCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer)
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(&hcloud.Server{
			ID:   123,
			Name: "test",
			ServerType: &hcloud.ServerType{
				ID:          45,
				Name:        "cax11",
				Description: "CAX11",
				Cores:       2,
				CPUType:     hcloud.CPUTypeShared,
				Memory:      4.0,
				Disk:        40,
				StorageType: hcloud.StorageTypeLocal,
			},
			Image: &hcloud.Image{
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
			},
			Datacenter: &hcloud.Datacenter{
				ID:   4,
				Name: "hel1-dc2",
				Location: &hcloud.Location{
					ID:          3,
					Name:        "hel1",
					Description: "Helsinki DC Park 1",
					NetworkZone: "eu-central",
					Country:     "FI",
					City:        "Helsinki",
					Latitude:    60.169855,
					Longitude:   24.938379,
				},
				Description: "Helsinki 1 virtual DC 2",
			},
			IncludedTraffic: 20 * util.Tebibyte,
			Protection:      hcloud.ServerProtection{Delete: true, Rebuild: true},
		}, nil, nil)

	out, err := fx.Run(cmd, []string{"test"})

	expOut := `ID:		123
Name:		test
Status:		
Created:	Mon Jan  1 00:00:00 UTC 0001 (a long while ago)
Server Type:	cax11 (ID: 45)
  ID:		45
  Name:		cax11
  Description:	CAX11
  Cores:	2
  CPU Type:	shared
  Memory:	4 GB
  Disk:		0 GB
  Storage Type:	local
  Public Net:
  IPv4:
    No Primary IPv4
  IPv6:
    No Primary IPv6
  Floating IPs:
    No Floating IPs
Private Net:
    No Private Networks
Volumes:
  No Volumes
Image:
  ID:		123
  Type:		system
  Status:	available
  Name:		test
  Description:	Test image
  Image size:	20.00 GB
  Disk size:	20 GB
  Created:	Fri Oct  6 12:00:00 UTC 1905 (a long while ago)
  OS flavor:	
  OS version:	-
  Rapid deploy:	no
Datacenter:
  ID:		4
  Name:		hel1-dc2
  Description:	Helsinki 1 virtual DC 2
  Location:
    Name:		hel1
    Description:	Helsinki DC Park 1
    Country:		FI
    City:		Helsinki
    Latitude:		60.169855
    Longitude:		24.938379
Traffic:
  Outgoing:	0 B
  Ingoing:	0 B
  Included:	20 TiB
Backup Window:	Backups disabled
Rescue System:	disabled
ISO:
  No ISO attached
Protection:
  Delete:	yes
  Rebuild:	yes
Labels:
  No labels
Placement Group:
  No Placement Group set
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
