package all_test

import (
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/all"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/all_list_response.json
var allListResponseJSON string

func TestListAll(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := all.ListCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.ServerListOpts{}).
		Return([]*hcloud.Server{
			{
				ID:     123,
				Name:   "my server",
				Status: hcloud.ServerStatusRunning,
				PublicNet: hcloud.ServerPublicNet{
					IPv4: hcloud.ServerPublicNetIPv4{
						IP: net.ParseIP("192.0.2.1"),
					},
				},
				Created:    time.Now().Add(-72 * time.Hour),
				Datacenter: &hcloud.Datacenter{Name: "hel1-dc2"},
			},
		}, nil)
	fx.Client.NetworkClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.NetworkListOpts{}).
		Return([]*hcloud.Network{
			{
				ID:      123,
				Name:    "test-net",
				IPRange: &net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: net.CIDRMask(24, 32)},
				Servers: []*hcloud.Server{{ID: 3421}},
				Created: time.Now().Add(-10 * time.Second),
			},
		}, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.LoadBalancerListOpts{}).
		Return([]*hcloud.LoadBalancer{
			{
				ID:               123,
				LoadBalancerType: &hcloud.LoadBalancerType{Name: "lb11"},
				Location:         &hcloud.Location{Name: "fsn1", NetworkZone: hcloud.NetworkZoneEUCentral},
				Name:             "foobar",
				PublicNet: hcloud.LoadBalancerPublicNet{
					IPv4: hcloud.LoadBalancerPublicNetIPv4{
						IP: net.ParseIP("192.0.2.1"),
					},
					IPv6: hcloud.LoadBalancerPublicNetIPv6{
						IP: net.IPv6zero,
					},
				},
				Created: time.Now().Add(-5 * time.Hour),
			},
		}, nil)
	fx.Client.CertificateClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.CertificateListOpts{}).
		Return([]*hcloud.Certificate{
			{
				ID:            123,
				Name:          "test",
				Type:          hcloud.CertificateTypeManaged,
				DomainNames:   []string{"example.com"},
				NotValidAfter: time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
				Created:       time.Now().Add(-10 * time.Hour),
			},
		}, nil)
	fx.Client.FirewallClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.FirewallListOpts{}).
		Return([]*hcloud.Firewall{
			{
				ID:        123,
				Name:      "test",
				Created:   time.Now().Add(-7 * time.Minute),
				Rules:     make([]hcloud.FirewallRule, 5),
				AppliedTo: make([]hcloud.FirewallResource, 2),
			},
		}, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.PrimaryIPListOpts{}).
		Return([]*hcloud.PrimaryIP{
			{
				ID:         123,
				Name:       "test",
				Created:    time.Now().Add(-2 * time.Hour),
				Datacenter: &hcloud.Datacenter{Name: "fsn1-dc14"},
				Type:       hcloud.PrimaryIPTypeIPv4,
				IP:         net.ParseIP("127.0.0.1"),
			},
		}, nil)
	fx.Client.FloatingIPClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.FloatingIPListOpts{}).
		Return([]*hcloud.FloatingIP{
			{
				ID:           123,
				Name:         "test",
				Created:      time.Now().Add(-24 * time.Hour),
				Type:         hcloud.FloatingIPTypeIPv4,
				IP:           net.ParseIP("127.0.0.1"),
				HomeLocation: &hcloud.Location{Name: "fsn1"},
			},
		}, nil)
	fx.Client.ImageClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.ImageListOpts{Type: []hcloud.ImageType{hcloud.ImageTypeBackup, hcloud.ImageTypeSnapshot}, IncludeDeprecated: true}).
		Return([]*hcloud.Image{
			{
				ID:           1,
				Type:         hcloud.ImageTypeBackup,
				Name:         "test",
				Created:      time.Date(2036, 8, 20, 12, 0, 0, 0, time.UTC),
				Architecture: hcloud.ArchitectureARM,
				DiskSize:     20,
			},
		}, nil)
	fx.Client.VolumeClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.VolumeListOpts{}).
		Return([]*hcloud.Volume{
			{
				ID:       123,
				Name:     "test",
				Location: &hcloud.Location{Name: "fsn1"},
				Size:     50,
				Created:  time.Now().Add(-500 * time.Hour),
			},
		}, nil)
	fx.Client.ISOClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.ISOListOpts{}).
		Return([]*hcloud.ISO{
			{
				ID:           123,
				Name:         "test",
				Type:         hcloud.ISOTypePrivate,
				Architecture: hcloud.Ptr(hcloud.ArchitectureARM),
			},
		}, nil)
	fx.Client.PlacementGroupClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.PlacementGroupListOpts{}).
		Return([]*hcloud.PlacementGroup{
			{
				ID:      123,
				Name:    "test",
				Created: time.Now().Add(-10 * time.Hour),
				Type:    hcloud.PlacementGroupTypeSpread,
				Servers: make([]int64, 5),
			},
		}, nil)
	fx.Client.SSHKeyClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.SSHKeyListOpts{}).
		Return([]*hcloud.SSHKey{
			{
				ID:      123,
				Name:    "test",
				Created: time.Now().Add(-2 * time.Hour),
			},
		}, nil)
	fx.Client.StorageBoxClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.StorageBoxListOpts{}).
		Return([]*hcloud.StorageBox{
			{
				ID:       123,
				Username: "u12345",
				Status:   hcloud.StorageBoxStatusActive,
				Name:     "test",
				Location: &hcloud.Location{Name: "fsn1"},
				Server:   "u1337.your-storagebox.de",
				System:   "FSN1-BX355",
				StorageBoxType: &hcloud.StorageBoxType{
					Name: "bx11",
				},
				Stats: hcloud.StorageBoxStats{
					Size: 42 * util.Gibibyte,
				},
				Labels: map[string]string{
					"environment":    "prod",
					"example.com/my": "label",
					"just-a-key":     "",
				},
				Protection: hcloud.StorageBoxProtection{
					Delete: false,
				},
				SnapshotPlan: &hcloud.StorageBoxSnapshotPlan{
					MaxSnapshots: 10,
					Minute:       1,
					Hour:         2,
					DayOfWeek:    hcloud.Ptr(time.Sunday),
					DayOfMonth:   hcloud.Ptr(4),
				},
				Created: time.Now().Add(-3 * time.Hour),
			},
		}, nil)
	fx.Client.ZoneClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.ZoneListOpts{}).
		Return([]*hcloud.Zone{
			{
				ID:      42,
				Name:    "example.com",
				Mode:    hcloud.ZoneModePrimary,
				Created: time.Now().Add(-1337 * time.Hour),
				Labels: map[string]string{
					"environment":    "prod",
					"example.com/my": "label",
					"just-a-key":     "",
				},
				Protection: hcloud.ZoneProtection{
					Delete: false,
				},
				RecordCount: 4,
				Registrar:   hcloud.ZoneRegistrarHetzner,
				Status:      hcloud.ZoneStatusOk,
				TTL:         10800,
			},
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{})

	expOut := `SERVERS
---
ID    NAME        STATUS    IPV4        IPV6   PRIVATE NET   DATACENTER   AGE
123   my server   running   192.0.2.1   -      -             hel1-dc2     3d 

IMAGES
---
ID   TYPE     NAME   DESCRIPTION   ARCHITECTURE   IMAGE SIZE   DISK SIZE   CREATED                   DEPRECATED
1    backup   test   -             arm            -            20 GB       2036-08-20 12:00:00 UTC   -         

PLACEMENT GROUPS
---
ID    NAME   SERVERS     TYPE     AGE
123   test   5 servers   spread   10h

PRIMARY IPS
---
ID    TYPE   NAME   IP          ASSIGNEE   DNS   AUTO DELETE   AGE
123   ipv4   test   127.0.0.1   -          -     no            2h 

ISOS
---
ID    NAME   DESCRIPTION   TYPE      ARCHITECTURE
123   test   -             private   arm         

VOLUMES
---
ID    NAME   SIZE    SERVER   LOCATION   AGE
123   test   50 GB   -        fsn1       20d

LOAD BALANCER
---
ID    NAME     HEALTH    IPV4        IPV6   TYPE   LOCATION   NETWORK ZONE   AGE
123   foobar   healthy   192.0.2.1   ::     lb11   fsn1       eu-central     5h 

FLOATING IPS
---
ID    TYPE   NAME   DESCRIPTION   IP          HOME   SERVER   DNS   AGE
123   ipv4   test   -             127.0.0.1   fsn1   -        -     1d 

NETWORKS
---
ID    NAME       IP RANGE       SERVERS    AGE
123   test-net   192.0.2.1/24   1 server   10s

FIREWALLS
---
ID    NAME   RULES COUNT   APPLIED TO COUNT                AGE
123   test   5 Rules       2 Servers | 0 Label Selectors   7m 

CERTIFICATES
---
ID    NAME   TYPE      DOMAIN NAMES   NOT VALID AFTER           AGE
123   test   managed   example.com    2036-08-20 12:00:00 UTC   10h

SSH KEYS
---
ID    NAME   FINGERPRINT   AGE
123   test   -             2h 

STORAGE BOXES
---
ID    NAME   USERNAME   SERVER                     TYPE   SIZE     LOCATION   AGE
123   test   u12345     u1337.your-storagebox.de   bx11   42 GiB   fsn1       3h 

ZONES
---
ID   NAME          STATUS   MODE      RECORD COUNT   AGE
42   example.com   ok       primary   4              55d

`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestListAllPaidJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := all.ListCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.ServerListOpts{}).
		Return([]*hcloud.Server{}, nil)
	fx.Client.LoadBalancerClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.LoadBalancerListOpts{}).
		Return([]*hcloud.LoadBalancer{}, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.PrimaryIPListOpts{}).
		Return([]*hcloud.PrimaryIP{}, nil)
	fx.Client.FloatingIPClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.FloatingIPListOpts{}).
		Return([]*hcloud.FloatingIP{}, nil)
	fx.Client.ImageClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.ImageListOpts{Type: []hcloud.ImageType{hcloud.ImageTypeBackup, hcloud.ImageTypeSnapshot}, IncludeDeprecated: true}).
		Return([]*hcloud.Image{
			{
				ID:           114690387,
				Name:         "debian-12",
				Description:  "Debian 12",
				Type:         hcloud.ImageTypeSystem,
				Status:       hcloud.ImageStatusAvailable,
				RapidDeploy:  true,
				OSVersion:    "12",
				OSFlavor:     "debian",
				DiskSize:     5,
				Architecture: hcloud.ArchitectureX86,
				Created:      time.Date(2023, 6, 13, 6, 0, 0, 0, time.UTC),
			},
		}, nil)
	fx.Client.VolumeClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.VolumeListOpts{}).
		Return([]*hcloud.Volume{}, nil)
	fx.Client.StorageBoxClient.EXPECT().
		AllWithOpts(gomock.Any(), hcloud.StorageBoxListOpts{}).
		Return([]*hcloud.StorageBox{
			{
				ID:       123,
				Username: "u12345",
				Status:   hcloud.StorageBoxStatusActive,
				Name:     "test",
				Location: &hcloud.Location{Name: "fsn1"},
				Server:   "u1337.your-storagebox.de",
				System:   "FSN1-BX355",
				StorageBoxType: &hcloud.StorageBoxType{
					Name: "bx11",
				},
				Stats: hcloud.StorageBoxStats{
					Size: 42 * util.Gibibyte,
				},
				Labels: map[string]string{
					"environment":    "prod",
					"example.com/my": "label",
					"just-a-key":     "",
				},
				Protection: hcloud.StorageBoxProtection{
					Delete: false,
				},
				SnapshotPlan: &hcloud.StorageBoxSnapshotPlan{
					MaxSnapshots: 10,
					Minute:       1,
					Hour:         2,
					DayOfWeek:    hcloud.Ptr(time.Sunday),
					DayOfMonth:   hcloud.Ptr(4),
				},
				Created: time.Date(2024, 4, 2, 10, 0, 0, 0, time.UTC),
			},
		}, nil)

	jsonOut, errOut, err := fx.Run(cmd, []string{"--paid", "-o=json"})

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.JSONEq(t, allListResponseJSON, jsonOut)
}
