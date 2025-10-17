package subaccount

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
)

func SuggestSubaccounts(client hcapi2.Client) cobra.CompletionFunc {
	return cmpl.SuggestCandidatesCtx(func(cmd *cobra.Command, args []string) []string {
		if len(args) == 0 {
			return nil
		}

		storageBox, _, err := client.StorageBox().Get(cmd.Context(), args[0])
		if err != nil || storageBox == nil {
			return nil
		}

		subaccounts, err := client.StorageBox().AllSubaccounts(context.Background(), storageBox)
		if err != nil {
			return nil
		}

		subaccountUsernames := make([]string, 0, len(subaccounts))
		for _, subaccount := range subaccounts {
			subaccountUsernames = append(subaccountUsernames, subaccount.Username)
		}
		return subaccountUsernames
	})
}
