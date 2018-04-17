package cli

import (
	"fmt"
	"strings"
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

var outputDescription = `Output can be controlled by the -o flag. Use -o noheader to suppress the table header.
Displayed columns and their order can be set with -o columns=%s see available columns below.`

func listLongDescription(intro string, columns []string) string {
	var colExample []string
	if len(columns) > 2 {
		colExample = columns[0:2]
	} else {
		colExample = columns
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\nColumns:\n - %s",
		intro,
		fmt.Sprintf(outputDescription, strings.Join(colExample, ",")),
		strings.Join(columns, "\n - "),
	)
}

func addListOutputFlag(cmd *cobra.Command, columns []string) {
	if len(columns) > 2 {
		columns = columns[0:2]
	}
	cmd.Flags().StringArrayP(
		"output",
		"o",
		[]string{},
		fmt.Sprintf(
			"output options: noheader|columns=%s,...",
			strings.Join(columns, ","),
		),
	)
}
