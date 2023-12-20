package util

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func YesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func NA(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func Datetime(t time.Time) string {
	return t.Local().Format(time.UnixDate)
}

func Age(t, currentTime time.Time) string {
	diff := currentTime.Sub(t)

	if int(diff.Hours()) >= 24 {
		days := int(diff.Hours()) / 24
		return fmt.Sprintf("%dd", days)
	}

	if int(diff.Hours()) > 0 {
		return fmt.Sprintf("%dh", int(diff.Hours()))
	}

	if int(diff.Minutes()) > 0 {
		return fmt.Sprintf("%dm", int(diff.Minutes()))
	}

	if int(diff.Seconds()) > 0 {
		return fmt.Sprintf("%ds", int(diff.Seconds()))
	}

	return "just now"
}

func ChainRunE(fns ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
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

func ExactlyOneSet(s string, ss ...string) bool {
	set := s != ""
	for _, s := range ss {
		if set && s != "" {
			return false
		}
		set = set || s != ""
	}
	return set
}

var outputDescription = `Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=%s (see available columns below).`

func ListLongDescription(intro string, columns []string) string {
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

func SplitLabel(label string) []string {
	return strings.SplitN(label, "=", 2)
}

// SplitLabelVars splits up label into key and value and returns them as separate return values.
// If label doesn't contain the `=` separator, SplitLabelVars returns the original string as key,
// with an empty value.
func SplitLabelVars(label string) (string, string) {
	parts := strings.SplitN(label, "=", 2)
	if len(parts) != 2 {
		return label, ""
	}
	return parts[0], parts[1]
}

func LabelsToString(labels map[string]string) string {
	var labelsString []string
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := labels[key]
		if value == "" {
			labelsString = append(labelsString, key)
		} else {
			labelsString = append(labelsString, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return strings.Join(labelsString, ", ")
}

// PrefixLines will prefix all individual lines in the text with the passed prefix.
func PrefixLines(text, prefix string) string {
	var lines []string

	for _, line := range strings.Split(text, "\n") {
		lines = append(lines, prefix+line)
	}

	return strings.Join(lines, "\n")
}

func DescribeFormat(object interface{}, format string) error {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	t, err := template.New("").Parse(format)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, object)
}

func DescribeJSON(object interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(object)
}

func DescribeYAML(object interface{}) error {
	enc := yaml.NewEncoder(os.Stdout)
	return enc.Encode(object)
}

// Wrap wraps the passed value in a map with the passed key.
//
// This is useful when working with JSON objects.
func Wrap(key string, value any) map[string]any {
	return map[string]any{key: value}
}

// ValidateRequiredFlags ensures that flags has values for all flags with
// the passed names.
//
// This function duplicates the functionality cobra provides when calling
// MarkFlagRequired. However, in some cases a flag cannot be marked as required
// in cobra, for example when it depends on other flags. In those cases this
// function comes in handy.
func ValidateRequiredFlags(flags *pflag.FlagSet, names ...string) error {
	var missingFlags []string

	for _, name := range names {
		if !flags.Changed(name) {
			missingFlags = append(missingFlags, `"`+name+`"`)
		}
	}
	if len(missingFlags) > 0 {
		return fmt.Errorf("hcloud: required flag(s) %s not set", strings.Join(missingFlags, ", "))
	}
	return nil
}
