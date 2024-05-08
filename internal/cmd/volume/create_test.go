package volume_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := volume.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.VolumeClient.EXPECT().
		Create(gomock.Any(), hcloud.VolumeCreateOpts{
			Name:     "test",
			Size:     20,
			Location: &hcloud.Location{Name: "fsn1"},
			Labels:   make(map[string]string),
		}).
		Return(hcloud.VolumeCreateResult{
			Volume: &hcloud.Volume{
				ID:       123,
				Name:     "test",
				Size:     20,
				Location: &hcloud.Location{Name: "fsn1"},
			},
			Action:      &hcloud.Action{ID: 321},
			NextActions: []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}, {ID: 1}, {ID: 2}, {ID: 3}})

	out, errOut, err := fx.Run(cmd, []string{"--name", "test", "--size", "20", "--location", "fsn1"})

	expOut := "Volume 123 created\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := volume.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.VolumeClient.EXPECT().
		Create(gomock.Any(), hcloud.VolumeCreateOpts{
			Name:     "test",
			Size:     20,
			Location: &hcloud.Location{Name: "fsn1"},
			Labels:   make(map[string]string),
		}).
		Return(hcloud.VolumeCreateResult{
			Volume: &hcloud.Volume{
				ID:       123,
				Name:     "test",
				Size:     20,
				Location: &hcloud.Location{Name: "fsn1"},
				Labels:   make(map[string]string),
				Created:  time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
				Status:   hcloud.VolumeStatusAvailable,
				Protection: hcloud.VolumeProtection{
					Delete: true,
				},
				Server: &hcloud.Server{ID: 123},
			},
			Action:      &hcloud.Action{ID: 321},
			NextActions: []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}, {ID: 1}, {ID: 2}, {ID: 3}})

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test", "--size", "20", "--location", "fsn1"})

	expOut := "Volume 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}

func TestCreateProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := volume.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	v := &hcloud.Volume{
		ID:       123,
		Name:     "test",
		Size:     20,
		Location: &hcloud.Location{Name: "fsn1"},
	}

	fx.Client.VolumeClient.EXPECT().
		Create(gomock.Any(), hcloud.VolumeCreateOpts{
			Name:     "test",
			Size:     20,
			Location: &hcloud.Location{Name: "fsn1"},
			Labels:   make(map[string]string),
		}).
		Return(hcloud.VolumeCreateResult{
			Volume:      v,
			Action:      &hcloud.Action{ID: 321},
			NextActions: []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), []*hcloud.Action{{ID: 321}, {ID: 1}, {ID: 2}, {ID: 3}})
	fx.Client.VolumeClient.EXPECT().
		ChangeProtection(gomock.Any(), v, hcloud.VolumeChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123})

	out, errOut, err := fx.Run(cmd, []string{"--name", "test", "--size", "20", "--location", "fsn1", "--enable-protection", "delete"})

	expOut := `Volume 123 created
Resource protection enabled for volume 123
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
