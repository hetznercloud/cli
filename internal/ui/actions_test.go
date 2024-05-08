package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestMessages(t *testing.T) {
	testCases := []struct {
		name          string
		action        *hcloud.Action
		wantAction    string
		wantResources string
	}{
		{
			name: "create_server",
			action: &hcloud.Action{
				ID:       1564532131,
				Command:  "create_server",
				Status:   hcloud.ActionStatusRunning,
				Progress: 0,
				Resources: []*hcloud.ActionResource{
					{ID: 46830545, Type: hcloud.ActionResourceTypeServer},
				},
			},
			wantAction:    "Waiting for create_server to complete",
			wantResources: "(server: 46830545)",
		},
		{
			name: "attach_volume",
			action: &hcloud.Action{
				ID:       1564532131,
				Command:  "attach_volume",
				Status:   hcloud.ActionStatusRunning,
				Progress: 0,
				Resources: []*hcloud.ActionResource{
					{ID: 46830545, Type: hcloud.ActionResourceTypeServer},
					{ID: 46830546, Type: hcloud.ActionResourceTypeVolume},
				},
			},
			wantAction:    "Waiting for attach_volume to complete",
			wantResources: "(server: 46830545, volume: 46830546)",
		},
		{
			name: "no resources",
			action: &hcloud.Action{
				ID:       1564532131,
				Command:  "create_server",
				Status:   hcloud.ActionStatusRunning,
				Progress: 0,
			},
			wantAction:    "Waiting for create_server to complete",
			wantResources: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actionMessage := ActionMessage(testCase.action)
			actionResourcesMessage := ActionResourcesMessage(testCase.action.Resources...)

			assert.Equal(t, testCase.wantAction, actionMessage)
			assert.Equal(t, testCase.wantResources, actionResourcesMessage)
		})
	}
}
