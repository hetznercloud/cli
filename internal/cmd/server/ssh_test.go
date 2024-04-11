package server_test

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestSSH(t *testing.T) {
	srv := hcloud.Server{
		ID:     42,
		Name:   "server1",
		Status: hcloud.ServerStatusRunning,
		PublicNet: hcloud.ServerPublicNet{
			IPv4: hcloud.ServerPublicNetIPv4{
				IP: net.ParseIP("192.168.0.2"),
			},
		},
	}

	preRun := func(t *testing.T, fx *testutil.Fixture) {
		fx.Client.ServerClient.EXPECT().
			Get(gomock.Any(), srv.Name).
			Return(&srv, nil, nil)

		fx.Config.EXPECT().SSHPath().Return("echo")
	}

	testutil.TestCommand(t, &server.SSHCmd, map[string]testutil.TestCase{
		"single arg": {
			Args:   []string{"ssh", srv.Name},
			PreRun: preRun,
			ExpOut: "-l root -p 22 192.168.0.2\n",
		},
		"many args": {
			Args:   []string{"ssh", srv.Name, "ls", "-al"},
			PreRun: preRun,
			ExpOut: "-l root -p 22 192.168.0.2 ls -al\n",
		},
	})

}
