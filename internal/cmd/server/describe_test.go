package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := server.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{
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
			Created:      time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
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
		Created:         time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(srv, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID:		123
Name:		test
Status:		
Created:	%s (%s)
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
  Created:	%s (%s)
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
`,
		util.Datetime(srv.Created), humanize.Time(srv.Created),
		util.Datetime(srv.Image.Created), humanize.Time(srv.Image.Created))

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
