package util

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"reflect"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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

func GrossPrice(price hcloud.Price) string {
	currencyDisplay := price.Currency

	// Currency is the ISO 4217 code, but we want to the show currency symbol
	switch price.Currency {
	case "EUR":
		currencyDisplay = "€"
	case "USD":
		currencyDisplay = "$"
	default:
		// unchanged
	}

	// The code/symbol and the amount are separated by a non-breaking space:
	// https://style-guide.europa.eu/en/content/-/isg/topic?identifier=7.3.3-rules-for-expressing-monetary-units#id370303__id370303_PositionISO
	return fmt.Sprintf("%s\u00a0%s", currencyDisplay, price.Gross)
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

func AnySet(ss ...string) bool {
	for _, s := range ss {
		if s != "" {
			return true
		}
	}
	return false
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
	for key, value := range IterateInOrder(labels) {
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
	var tail string
	if len(text) > 0 && text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
		tail = "\n"
	}

	return prefix + strings.ReplaceAll(text, "\n", "\n"+prefix) + tail
}

func DescribeFormat(w io.Writer, object interface{}, format string) error {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	t, err := template.New("").Parse(format)
	if err != nil {
		return err
	}
	return t.Execute(w, object)
}

func DescribeJSON(w io.Writer, object interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(object)
}

func DescribeYAML(w io.Writer, object interface{}) error {
	enc := yaml.NewEncoder(w)
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

// AddGroup adds a group to the passed command and adds the passed commands to the group.
func AddGroup(cmd *cobra.Command, id string, title string, groupCmds ...*cobra.Command) {
	if !strings.HasSuffix(title, ":") {
		title += ":"
	}
	cmd.AddGroup(&cobra.Group{ID: id, Title: title})
	for _, groupCmd := range groupCmds {
		groupCmd.GroupID = id
	}
	cmd.AddCommand(groupCmds...)
}

// ToKebabCase converts the passed string to kebab-case.
func ToKebabCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

func IsNil(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

func FilterNil[T any](values []T) []T {
	var filtered []T
	for _, v := range values {
		if !IsNil(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// SliceDiff returns the difference between the two passed slices. The returned slice contains all elements that
// are present in a but not in b. The order of a is preserved.
func SliceDiff[S ~[]E, E cmp.Ordered](a, b []E) []E {
	var diff S
	seen := make(map[E]struct{})
	for _, v := range b {
		seen[v] = struct{}{}
	}
	for _, v := range a {
		if _, ok := seen[v]; ok {
			continue
		}
		diff = append(diff, v)
	}
	return diff
}

// RemoveDuplicates removes duplicates from the passed slice while preserving the order of the elements.
// The first occurrence of an element is kept, all following occurrences are removed.
func RemoveDuplicates[S ~[]E, E cmp.Ordered](values S) S {
	var unique S
	seen := make(map[E]struct{})
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		unique = append(unique, v)
	}
	return unique
}

func AnyToAnySlice(a any) []any {
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Slice {
		return nil
	}
	s := make([]any, val.Len())
	for i := 0; i < val.Len(); i++ {
		s[i] = val.Index(i).Interface()
	}
	return s
}

func AnyToStringSlice(a any) []string {
	var s []string
	for _, v := range AnyToAnySlice(a) {
		s = append(s, fmt.Sprint(v))
	}
	return s
}

func ToStringSlice(a []any) []string {
	var s []string
	for _, v := range a {
		s = append(s, fmt.Sprint(v))
	}
	return s
}

func ToAnySlice[T any](a []T) []any {
	var s []any
	for _, v := range a {
		s = append(s, any(v))
	}
	return s
}

// ParseBoolLenient parses the passed string as a boolean. It is different from strconv.ParseBool in that it
// is case-insensitive and also accepts "yes"/"y" and "no"/"n" as valid values.
func ParseBoolLenient(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "t", "yes", "y", "1":
		return true, nil
	case "false", "f", "no", "n", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}

// ToBoolE converts the provided value to a bool. It is more lenient than [cast.ToBoolE] (see [ParseBoolLenient]).
func ToBoolE(val any) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		return ParseBoolLenient(v)
	default:
		return cast.ToBoolE(val)
	}
}

// ToStringSliceDelimited is like [AnyToStringSlice] but also accepts a string that is comma-separated.
func ToStringSliceDelimited(val any) []string {
	switch v := val.(type) {
	case []string:
		return v
	case string:
		return strings.Split(v, ",")
	default:
		return AnyToStringSlice(val)
	}
}

// IterateInOrder returns an iterator that iterates over the map in order of the keys.
func IterateInOrder[M ~map[K]V, K cmp.Ordered, V any](m M) iter.Seq2[K, V] {
	keys := slices.Sorted(maps.Keys(m))
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

func FormatHcloudError(err error) string {
	var hcloudErr hcloud.Error
	if !errors.As(err, &hcloudErr) {
		return err.Error()
	}

	switch details := hcloudErr.Details.(type) {
	case hcloud.ErrorDetailsInvalidInput:
		var errBuilder strings.Builder
		errBuilder.WriteString(hcloudErr.Error())
		for _, field := range details.Fields {
			fieldMsg := strings.Join(field.Messages, ", ")
			errBuilder.WriteString(fmt.Sprintf("\n- %s: %s", field.Name, fieldMsg))
		}
		return errBuilder.String()

	default:
		return err.Error()
	}
}

func OptionalString(s *string, defaultValue string) string {
	if s == nil || *s == "" {
		return defaultValue
	}
	return *s
}
