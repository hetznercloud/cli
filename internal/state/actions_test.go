package state

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

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
		DoAndReturn(func(_ context.Context, handleUpdate func(update *hcloud.Action) error, _ ...*hcloud.Action) error {
			require.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusRunning
			require.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusSuccess
			require.NoError(t, handleUpdate(action))

			return nil
		})

	var stderr strings.Builder
	require.NoError(t, waitForActions(t.Context(), client, &stderr, action))

	assert.Equal(t,
		strings.Join([]string{
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ...\n",
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ... done\n",
		}, ""),
		stderr.String(),
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
		DoAndReturn(func(_ context.Context, handleUpdate func(update *hcloud.Action) error, _ ...*hcloud.Action) error {
			require.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusRunning
			require.NoError(t, handleUpdate(action))
			action.Status = hcloud.ActionStatusError
			action.ErrorCode = "action_failed"
			action.ErrorMessage = "action failed"
			require.NoError(t, handleUpdate(action))

			return nil
		})

	var stderr strings.Builder
	err := waitForActions(t.Context(), client, &stderr, action)
	var actionWaitErr *ActionWaitError
	require.ErrorAs(t, err, &actionWaitErr)
	require.Len(t, actionWaitErr.Failures, 1)
	assert.Equal(t, action.ID, actionWaitErr.Failures[0].ActionID)
	require.ErrorContains(t, actionWaitErr.Failures[0].Err, "action failed")

	assert.Equal(t,
		strings.Join([]string{
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ...\n",
			"Waiting for attach_volume (server: 46830545, volume: 46830546) ... failed\n",
		}, ""),
		stderr.String(),
	)
}

func TestCollectActionFailuresWaitsForEveryAction(t *testing.T) {
	ctrl := gomock.NewController(t)

	actions := []*hcloud.Action{{ID: 1}, {ID: 2}}
	client := hcapi2_mock.NewMockActionClient(ctrl)
	client.EXPECT().
		WaitForFunc(gomock.Any(), gomock.Any(), actions[0], actions[1]).
		DoAndReturn(func(_ context.Context, handleUpdate func(*hcloud.Action) error, _ ...*hcloud.Action) error {
			for _, action := range actions {
				action.Status = hcloud.ActionStatusError
				action.ErrorCode = "action_failed"
				action.ErrorMessage = "failed"
				require.NoError(t, handleUpdate(action))
			}
			return nil
		})

	err := waitForActionsQuiet(t.Context(), client, actions...)
	var actionWaitErr *ActionWaitError
	require.ErrorAs(t, err, &actionWaitErr)
	assert.Len(t, actionWaitErr.Failures, 2)
	assert.NoError(t, actionWaitErr.Cause)
}
