package certificate

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var RetryCmd = base.Cmd{
	BaseCobraCommand: func(_ hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:                   "retry <certificate>",
			Short:                 "Retry a Managed Certificate's issuance",
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		certificate, _, err := s.Client().Certificate().Get(s, idOrName)
		if err != nil {
			return err
		}
		if certificate == nil {
			return fmt.Errorf("Certificate not found: %s", idOrName)
		}

		action, _, err := s.Client().Certificate().RetryIssuance(s, certificate)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Retried issuance of Certificate %s\n", certificate.Name)
		return nil
	},
}
