package server_test

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

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
	out, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	expOut := `Server 1234 created
IPv4: 192.0.2.1
`
	assert.Equal(t, expOut, out)
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

	server := &hcloud.Server{
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
		Return(server, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).Return(nil)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 234}}).Return(nil)

	fx.Client.ServerClient.EXPECT().
		ChangeProtection(gomock.Any(), server, hcloud.ServerChangeProtectionOpts{
			Rebuild: hcloud.Ptr(true), Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{
			ID: 1337,
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 1337}).Return(nil)

	fx.Client.ServerClient.EXPECT().
		EnableBackup(gomock.Any(), server, "").
		Return(&hcloud.Action{
			ID: 42,
		}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 42}).Return(nil)

	args := []string{"--name", "cli-test", "--type", "cx11", "--image", "ubuntu-20.04", "--enable-protection", "rebuild,delete", "--enable-backup"}
	out, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	expOut := `Server 1234 created
Resource protection enabled for server 1234
Backups enabled for server 1234
IPv4: 192.0.2.1
`
	assert.Equal(t, expOut, out)
}
