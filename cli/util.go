package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
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
			if fn == nil {
				continue
			}
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

var outputDescription = `Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=%s (see available columns below).`

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

func splitLabel(label string) []string {
	return strings.SplitN(label, "=", 2)
}

func labelsToString(labels map[string]string) string {
	var labelsString []string
	for key, value := range labels {
		if value == "" {
			labelsString = append(labelsString, key)
		} else {
			labelsString = append(labelsString, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return strings.Join(labelsString, ", ")
}

func describeFormat(object interface{}, format string) error {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	t, err := template.New("").Parse(format)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, object)
}

func describeJSON(object interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(object)
}
