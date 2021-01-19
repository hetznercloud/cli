package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := hcapi2.NewMockClient(ctrl)
	actionWaiter := state.NewMockActionWaiter(ctrl)
	tokenEnsurer := state.NewMockTokenEnsurer(ctrl)

	cmd := newCreateCommand(
		context.Background(),
		client,
		tokenEnsurer,
		actionWaiter,
	)

	tokenEnsurer.EXPECT().EnsureToken(gomock.Any(), gomock.Any()).Return(nil)
	client.ImageClient.EXPECT().
		Get(gomock.Any(), "ubuntu-20.04").
		Return(&hcloud.Image{}, nil, nil)
	client.ServerClient.EXPECT().
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
	actionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).Return(nil)
	actionWaiter.EXPECT().WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 234}}).Return(nil)

	args := []string{"--name", "cli-test", "--type", "cx11", "--image", "ubuntu-20.04"}
	cmd.SetArgs(args)

	out, err := captureStdout(func() error {
		return cmd.Execute()
	})

	assert.NoError(t, err)
	expOut := `Server 1234 created
IPv4: 192.0.2.1
`
	assert.Equal(t, expOut, out)
}

func captureStdout(fn func() error) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("capture stdout: %v", err)
	}

	origOut := os.Stdout
	defer func() {
		os.Stdout = origOut
	}()

	buf := &bytes.Buffer{}
	os.Stdout = w

	copyDone := make(chan struct{})
	var copyErr error
	go func() {
		defer close(copyDone)

		copied, err := io.Copy(buf, r)
		if err != nil {
			copyErr = fmt.Errorf("capture stdout: %v, copied: %d\n", err, copied)
			return
		}
	}()

	err = fn()

	if copyErr != nil {
		return "", copyErr
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("capture stdout close pipe reader: %v", err)
	}

	<-copyDone

	return buf.String(), err
}
