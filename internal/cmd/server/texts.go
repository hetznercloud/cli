package server

import (
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func warningDeprecatedServerType(serverType *hcloud.ServerType) string {
	if !serverType.IsDeprecated() {
		return ""
	}

	if time.Now().After(serverType.UnavailableAfter()) {
		return fmt.Sprintf("Attention: The Server Type %q is deprecated and can no longer be ordered. Existing servers of that plan will continue to work as before and no action is required on your part. It is possible to migrate this Server to another Server Type by using the \"hcloud server change-type\" command.\n\n", serverType.Name)
	}

	return fmt.Sprintf("Attention: The Server Type %q is deprecated and will no longer be available for order as of %s. Existing servers of that plan will continue to work as before and no action is required on your part. It is possible to migrate this Server to another Server Type by using the \"hcloud server change-type\" command.\n\n", serverType.Name, serverType.UnavailableAfter().Format(time.DateOnly))
}
