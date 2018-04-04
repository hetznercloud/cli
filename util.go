package cli

import (
	"time"

	"github.com/spf13/cobra"
)

func yesno(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func na(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func datetime(t time.Time) string {
	return t.Local().Format(time.UnixDate)
}

func chainRunE(fns ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range fns {
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
