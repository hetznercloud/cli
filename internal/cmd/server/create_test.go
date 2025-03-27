package server_test

import (
	"context"
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

//go:embed testdata/create_response.json
var createResponseJSON string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.CreateCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cx22").
		Return(&hcloud.ServerType{Architecture: hcloud.ArchitectureX86}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "ubuntu-20.04", hcloud.ArchitectureX86).
		Return(&hcloud.Image{}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, opts hcloud.ServerCreateOpts) {
			assert.Equal(t, "cli-test", opts.Name)
		}).
		Return(hcloud.ServerCreateResult{
			Server: &hcloud.Server{
				ID: 1234,
				PublicNet: hcloud.ServerPublicNet{
					IPv4: hcloud.ServerPublicNetIPv4FromSchema(schema.ServerPublicNetIPv4{
						IP: "192.0.2.1",
					}),
					IPv6: hcloud.ServerPublicNetIPv6FromSchema(schema.ServerPublicNetIPv6{
						IP: "2001:0db8:c013:4d58::/64",
					}),
				},
				PrivateNet: []hcloud.ServerPrivateNet{
					hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{
						Network: 4461841,
						IP:      "10.1.0.2",
					}),
					hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{
						Network: 4461842,
						IP:      "10.2.0.2",
					}),
				},
			},
			Action:       &hcloud.Action{ID: 123},
			NextActions:  []*hcloud.Action{{ID: 234}},
			RootPassword: "password",
		}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), int64(1234)).
		Return(&hcloud.Server{
			ID: 1234,
			PublicNet: hcloud.ServerPublicNet{
				IPv4: hcloud.ServerPublicNetIPv4FromSchema(schema.ServerPublicNetIPv4{
					IP: "192.0.2.1",
				}),
				IPv6: hcloud.ServerPublicNetIPv6FromSchema(schema.ServerPublicNetIPv6{
					IP: "2001:0db8:c013:4d58::/64",
				}),
			},
			PrivateNet: []hcloud.ServerPrivateNet{
				hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{
					Network: 4461841,
					IP:      "10.1.0.2",
				}),
				hcloud.ServerPrivateNetFromSchema(schema.ServerPrivateNet{
					Network: 4461842,
					IP:      "10.2.0.2",
				}),
			},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 234}}).
		Return(nil)
	fx.Client.NetworkClient.EXPECT().
		Name(int64(4461841)).
		Return("foo")
	fx.Client.NetworkClient.EXPECT().
		Name(int64(4461842)).
		Return("bar")

	args := []string{"--name", "cli-test", "--type", "cx22", "--image", "ubuntu-20.04"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	expOut := `Server 1234 created
IPv4: 192.0.2.1
IPv6: 2001:db8:c013:4d58::1
IPv6 Network: 2001:db8:c013:4d58::/64
Private Networks:
	- 10.1.0.2 (foo)
	- 10.2.0.2 (bar)
Root password: password
`
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := server.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{
		ID:   1234,
		Name: "cli-test",
		PublicNet: hcloud.ServerPublicNet{
			IPv4: hcloud.ServerPublicNetIPv4{
				IP: net.ParseIP("192.0.2.1"),
			},
		},
		Created: time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
		Labels:  make(map[string]string),
		Datacenter: &hcloud.Datacenter{
			ID:   1,
			Name: "fsn1-dc14",
			Location: &hcloud.Location{
				ID:   1,
				Name: "fsn1",
			},
		},
		ServerType: &hcloud.ServerType{
			ID:           1,
			Name:         "cx22",
			Cores:        1,
			CPUType:      "shared",
			Memory:       2,
			Disk:         20,
			StorageType:  "local",
			Architecture: hcloud.ArchitectureX86,
		},
		Image: &hcloud.Image{
			ID:          1,
			Type:        "system",
			Status:      "available",
			Name:        "ubuntu-20.04",
			Description: "Ubuntu 20.04",
			Labels:      make(map[string]string),
			OSFlavor:    "ubuntu",
			OSVersion:   "20.04",
			RapidDeploy: true,
			Protection: hcloud.ImageProtection{
				Delete: true,
			},
		},
		ISO: &hcloud.ISO{
			ID:          1,
			Name:        "FreeBSD-11.0-RELEASE-amd64-dvd1",
			Description: "FreeBSD 11.0 x64",
			Type:        "public",
		},
		RescueEnabled: true,
		Locked:        true,
		Status:        hcloud.ServerStatusRunning,
	}

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cx22").
		Return(&hcloud.ServerType{Architecture: hcloud.ArchitectureX86}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "ubuntu-20.04", hcloud.ArchitectureX86).
		Return(&hcloud.Image{}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, opts hcloud.ServerCreateOpts) {
			assert.Equal(t, "cli-test", opts.Name)
		}).
		Return(hcloud.ServerCreateResult{
			Server:       srv,
			RootPassword: "secret",
			Action:       &hcloud.Action{ID: 123},
			NextActions:  []*hcloud.Action{{ID: 234}},
		}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), int64(1234)).
		Return(srv, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 234}}).
		Return(nil)

	args := []string{"-o=json", "--name", "cli-test", "--type", "cx22", "--image", "ubuntu-20.04"}
	jsonOut, out, err := fx.Run(cmd, args)

	expOut := "Server 1234 created\n"

	require.NoError(t, err)
	assert.Equal(t, expOut, out)

	assert.JSONEq(t, createResponseJSON, jsonOut)
}

func TestCreateProtectionBackup(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.CreateCmd.CobraCommand(fx.State())

	fx.ExpectEnsureToken()

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cx22").
		Return(&hcloud.ServerType{Architecture: hcloud.ArchitectureX86}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "ubuntu-20.04", hcloud.ArchitectureX86).
		Return(&hcloud.Image{}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, opts hcloud.ServerCreateOpts) {
			assert.Equal(t, "cli-test", opts.Name)
		}).
		Return(hcloud.ServerCreateResult{
			Server: &hcloud.Server{
				ID: 1234,
				PublicNet: hcloud.ServerPublicNet{
					IPv4: hcloud.ServerPublicNetIPv4{
						IP: net.ParseIP("192.0.2.1"),
					},
				},
			},
			Action:      &hcloud.Action{ID: 123},
			NextActions: []*hcloud.Action{{ID: 234}},
		}, nil, nil)

	srv := &hcloud.Server{
		ID: 1234,
		PublicNet: hcloud.ServerPublicNet{
			IPv4: hcloud.ServerPublicNetIPv4{
				IP: net.ParseIP("192.0.2.1"),
			},
		},
		Protection: hcloud.ServerProtection{
			Delete:  true,
			Rebuild: true,
		},
	}

	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), int64(1234)).
		Return(srv, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 123}, {ID: 234}}).
		Return(nil)

	fx.Client.ServerClient.EXPECT().
		ChangeProtection(gomock.Any(), srv, hcloud.ServerChangeProtectionOpts{
			Rebuild: hcloud.Ptr(true), Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{
			ID: 1337,
		}, nil, nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 1337}).Return(nil)

	fx.Client.ServerClient.EXPECT().
		EnableBackup(gomock.Any(), srv, "").
		Return(&hcloud.Action{
			ID: 42,
		}, nil, nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 42}).Return(nil)

	args := []string{"--name", "cli-test", "--type", "cx22", "--image", "ubuntu-20.04", "--enable-protection", "rebuild,delete", "--enable-backup"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	expOut := `Server 1234 created
Resource protection enabled for Server 1234
Backups enabled for Server 1234
IPv4: 192.0.2.1
`
	assert.Equal(t, expOut, out)
}
