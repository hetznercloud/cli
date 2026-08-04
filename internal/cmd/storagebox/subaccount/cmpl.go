package subaccount

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
)

func SuggestSubaccounts(client hcapi2.Client) cobra.CompletionFunc {
	return cmpl.SuggestCandidatesCtxE(func(cmd *cobra.Command, args []string) ([]string, error) {
		if len(args) == 0 {
			return nil, nil
		}

		storageBox, _, err := client.StorageBox().Get(cmd.Context(), args[0])
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, nil
		}

		subaccounts, err := client.StorageBox().AllSubaccounts(cmd.Context(), storageBox)
		if err != nil {
			return nil, err
		}

		subaccountUsernames := make([]string, 0, len(subaccounts))
		for _, subaccount := range subaccounts {
			subaccountUsernames = append(subaccountUsernames, subaccount.Username)
		}
		return subaccountUsernames, nil
	})
}
