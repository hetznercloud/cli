package util

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
)

func TestMergeNextActions(t *testing.T) {
	action := &hcloud.Action{ID: 1}
	next_actions := []*hcloud.Action{{ID: 2}, {ID: 3}}

	actions := MergeNextActions(action, next_actions)

	assert.Equal(t, []*hcloud.Action{{ID: 1}, {ID: 2}, {ID: 3}}, actions)
}
