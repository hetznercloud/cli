package volume

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

//go:embed testdata/create_response.json
var createResponseJson string

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
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
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}})

	out, _, err := fx.Run(cmd, []string{"--name", "test", "--size", "20", "--location", "fsn1"})

	expOut := "Volume 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	response, err := testutil.MockResponse(&schema.VolumeCreateResponse{
		Volume: schema.Volume{
			ID:       123,
			Name:     "test",
			Size:     20,
			Location: schema.Location{Name: "fsn1"},
			Labels:   make(map[string]string),
			Created:  time.Date(2016, 1, 30, 23, 50, 0, 0, time.UTC),
			Status:   string(hcloud.VolumeStatusAvailable),
			Protection: schema.VolumeProtection{
				Delete: true,
			},
			Server: hcloud.Ptr(int64(123)),
		},
		Action: &schema.Action{
			ID: 321,
		},
		NextActions: make([]schema.Action, 0),
	})

	if err != nil {
		t.Fatal(err)
	}

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
		}, response, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}})

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test", "--size", "20", "--location", "fsn1"})

	expOut := "Volume 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, createResponseJson, jsonOut)
}

func TestCreateProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	volume := &hcloud.Volume{
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
			Volume:      volume,
			Action:      &hcloud.Action{ID: 321},
			NextActions: []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 321})
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}})
	fx.Client.VolumeClient.EXPECT().
		ChangeProtection(gomock.Any(), volume, hcloud.VolumeChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), &hcloud.Action{ID: 123})

	out, _, err := fx.Run(cmd, []string{"--name", "test", "--size", "20", "--location", "fsn1", "--enable-protection", "delete"})

	expOut := `Volume 123 created
Resource protection enabled for volume 123
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
