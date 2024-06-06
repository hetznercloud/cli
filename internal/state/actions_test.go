package state

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	hcapi2_mock "github.com/hetznercloud/cli/internal/hcapi2/mock"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestWaitForActionsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	action := &hcloud.Action{
		ID:       1564532131,
		Command:  "attach_volume",
		Status:   hcloud.ActionStatusRunning,
		Progress: 0,
		Resources: []*hcloud.ActionResource{
			{ID: 46830545, Type: hcloud.ActionResourceTypeServer},
			{ID: 46830546, Type: hcloud.ActionResourceTypeVolume},
		},
	}

	client := hcapi2_mock.NewMockActionClient(ctrl)

	client.EXPECT().
		WaitForFunc(gomock.Any(), gomock.Any(), action).
		DoAndReturn(func(ctx context.Context, handleUpdate func(update *hcloud.Action) error, actions ...*hcloud.Action) error {
			assert.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusRunning
			assert.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusSuccess
			assert.NoError(t, handleUpdate(action))

			return nil
		})

	stderr := captureStderr(t, func() {
		waitForActions(client, context.Background(), action)
	})

	assert.Equal(t,
		strings.Join([]string{
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ...\n",
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ... done\n",
		}, ""),
		stderr,
	)
}

func TestWaitForActionsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	action := &hcloud.Action{
		ID:       1564532131,
		Command:  "attach_volume",
		Status:   hcloud.ActionStatusRunning,
		Progress: 0,
		Resources: []*hcloud.ActionResource{
			{ID: 46830545, Type: hcloud.ActionResourceTypeServer},
			{ID: 46830546, Type: hcloud.ActionResourceTypeVolume},
		},
	}

	client := hcapi2_mock.NewMockActionClient(ctrl)
	client.EXPECT().
		WaitForFunc(gomock.Any(), gomock.Any(), action).
		DoAndReturn(func(ctx context.Context, handleUpdate func(update *hcloud.Action) error, actions ...*hcloud.Action) error {
			assert.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusRunning
			assert.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusError
			action.ErrorCode = "action_failed"
			action.ErrorMessage = "action failed"
			assert.Error(t, handleUpdate(action))

			return action.Error()
		})

	stderr := captureStderr(t, func() {
		waitForActions(client, context.Background(), action)
	})

	assert.Equal(t,
		strings.Join([]string{
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ...\n",
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ... failed\n",
		}, ""),
		stderr,
	)
}

func captureStderr(t *testing.T, next func()) string {
	t.Helper()

	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stderrOrig := os.Stderr
	os.Stderr = pipeWriter
	defer func() {
		os.Stderr = stderrOrig
	}()

	next()

	if err := pipeWriter.Close(); err != nil {
		t.Fatal(err)
	}

	stderrOutput, err := io.ReadAll(pipeReader)
	if err != nil {
		t.Fatal(err)
	}

	return string(stderrOutput)
}
