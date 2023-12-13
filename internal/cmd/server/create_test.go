package server_test

import (
	"context"
	_ "embed"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter,
	)

	fx.ExpectEnsureToken()

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cx11").
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
	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), int64(1234)).
		Return(&hcloud.Server{
			ID: 1234,
			PublicNet: hcloud.ServerPublicNet{
				IPv4: hcloud.ServerPublicNetIPv4{
					IP: net.ParseIP("192.0.2.1"),
				},
			},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).Return(nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 234}}).Return(nil)

	args := []string{"--name", "cli-test", "--type", "cx11", "--image", "ubuntu-20.04"}
	out, _, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	expOut := `Server 1234 created
IPv4: 192.0.2.1
`
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := server.CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter,
	)
	fx.ExpectEnsureToken()

	response, err := testutil.MockResponse(&schema.ServerCreateResponse{
		Server: schema.Server{
			ID:   1234,
			Name: "cli-test",
			PublicNet: schema.ServerPublicNet{
				IPv4: schema.ServerPublicNetIPv4{
					IP: "192.0.2.1",
				},
			},
			Created: time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
			Labels:  make(map[string]string),
			Datacenter: schema.Datacenter{
				ID:   1,
				Name: "fsn1-dc14",
				Location: schema.Location{
					ID:   1,
					Name: "fsn1",
				},
			},
			ServerType: schema.ServerType{
				ID:           1,
				Name:         "cx11",
				Cores:        1,
				CPUType:      "shared",
				Memory:       2,
				Disk:         20,
				StorageType:  "local",
				Architecture: string(hcloud.ArchitectureX86),
			},
			Image: &schema.Image{
				ID:          1,
				Type:        "system",
				Status:      "available",
				Name:        hcloud.Ptr("ubuntu-20.04"),
				Description: "Ubuntu 20.04",
				Deprecated:  nil,
				Labels:      make(map[string]string),
				OSFlavor:    "ubuntu",
				OSVersion:   hcloud.Ptr("20.04"),
				RapidDeploy: true,
				Protection: schema.ImageProtection{
					Delete: true,
				},
			},
			ISO: &schema.ISO{
				ID:          1,
				Name:        "FreeBSD-11.0-RELEASE-amd64-dvd1",
				Description: "FreeBSD 11.0 x64",
				Type:        "public",
				Deprecated:  nil,
			},
			RescueEnabled: true,
			Locked:        true,
			Status:        string(hcloud.ServerStatusRunning),
		},
		NextActions:  make([]schema.Action, 0),
		RootPassword: hcloud.Ptr("secret"),
		Action:       schema.Action{ID: 123},
	})

	if err != nil {
		t.Fatal(err)
	}

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cx11").
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
		}, response, nil)
	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), int64(1234)).
		Return(&hcloud.Server{
			ID: 1234,
			PublicNet: hcloud.ServerPublicNet{
				IPv4: hcloud.ServerPublicNetIPv4{
					IP: net.ParseIP("192.0.2.1"),
				},
			},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).Return(nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 234}}).Return(nil)

	args := []string{"-o=json", "--name", "cli-test", "--type", "cx11", "--image", "ubuntu-20.04"}
	jsonOut, out, err := fx.Run(cmd, args)

	expOut := "Server 1234 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}

func TestCreateProtectionBackup(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter,
	)

	fx.ExpectEnsureToken()

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cx11").
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
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).Return(nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 234}}).Return(nil)

	fx.Client.ServerClient.EXPECT().
		ChangeProtection(gomock.Any(), srv, hcloud.ServerChangeProtectionOpts{
			Rebuild: hcloud.Ptr(true), Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{
			ID: 1337,
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 1337}).Return(nil)

	fx.Client.ServerClient.EXPECT().
		EnableBackup(gomock.Any(), srv, "").
		Return(&hcloud.Action{
			ID: 42,
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 42}).Return(nil)

	args := []string{"--name", "cli-test", "--type", "cx11", "--image", "ubuntu-20.04", "--enable-protection", "rebuild,delete", "--enable-backup"}
	out, _, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	expOut := `Server 1234 created
Resource protection enabled for server 1234
Backups enabled for server 1234
IPv4: 192.0.2.1
`
	assert.Equal(t, expOut, out)
}
