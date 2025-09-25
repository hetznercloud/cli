package server

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/deprecationutil"
)

const ChangeDeprecatedServerTypeMessage = (`Existing servers of that plan will ` +
	`continue to work as before and no action is required on your part. ` +
	`It is possible to migrate this Server to another Server Type by using ` +
	`the "hcloud server change-type" command.`)

func deprecatedServerTypeWarning(serverType *hcloud.ServerType, locationName string) string {
	warnText, _ := deprecationutil.ServerTypeWarning(serverType, locationName)
	if warnText == "" {
		return ""
	}

	return fmt.Sprintf("Attention: %s %s\n\n", warnText, ChangeDeprecatedServerTypeMessage)
}
