package storagebox

import (
	"fmt"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func warningDeprecatedStorageBoxType(storageBoxType *hcloud.StorageBoxType) string {
	if !storageBoxType.IsDeprecated() {
		return ""
	}

	if time.Now().After(storageBoxType.UnavailableAfter()) {
		return fmt.Sprintf("Attention: The Storage Box Type %q is deprecated and can no longer be ordered. Existing Storage Boxes of that plan will continue to work as before and no action is required on your part. It is possible to migrate this Storage Box to another Storage Box Type by using the \"hcloud storage-box change-type\" command.\n\n", storageBoxType.Name)
	}

	return fmt.Sprintf("Attention: The Storage Box Type %q is deprecated and will no longer be available for order as of %s. Existing Storage Boxes of that plan will continue to work as before and no action is required on your part. It is possible to migrate this Storage Box to another Storage Box Type by using the \"hcloud storage-box change-type\" command.\n\n", storageBoxType.Name, storageBoxType.UnavailableAfter().Format(time.DateOnly))
}
