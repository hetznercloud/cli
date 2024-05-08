package util

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

func MergeNextActions(action *hcloud.Action, nextActions []*hcloud.Action) []*hcloud.Action {
	all := make([]*hcloud.Action, 0, 1+len(nextActions))
	all = append(all, action)
	all = append(all, nextActions...)
	return all
}
