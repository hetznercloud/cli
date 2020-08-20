package cli

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/fatih/structs"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/spf13/cobra"
)

const flagName = "output"

type outputOption struct {
	Name   string
	Values []string
}

func outputOptionNoHeader() outputOption {
	return outputOption{Name: "noheader"}
}

func outputOptionJSON() outputOption {
	return outputOption{Name: "json"}
}

func outputOptionFormat() outputOption {
	return outputOption{Name: "format"}
}

func outputOptionColumns(columns []string) outputOption {
	return outputOption{Name: "columns", Values: columns}
}

func addOutputFlag(cmd *cobra.Command, options ...outputOption) {
	var (
		names  []string
		values []string
	)
	for _, option := range options {
		name := option.Name
		if option.Values != nil {
			name += "=..."
			values = append(values, option.Values...)
		}
		names = append(names, name)
	}
	cmd.Flags().StringArrayP(
		flagName,
		"o",
		[]string{},
		fmt.Sprintf("output options: %s", strings.Join(names, "|")),
	)
	cmd.RegisterFlagCompletionFunc(flagName, cmpl.SuggestCandidates(values...))
	cmd.PreRunE = chainRunE(cmd.PreRunE, validateOutputFlag(options))
}

func validateOutputFlag(options []outputOption) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		validOptions := map[string]map[string]bool{}
		for _, option := range options {
			if option.Values == nil {
				// all values allowed
				validOptions[option.Name] = nil
			} else {
				// only discrete values allowed
				validOptions[option.Name] = map[string]bool{}
				for _, value := range option.Values {
					validOptions[option.Name][value] = true
				}
			}
		}

		flagValues, err := cmd.Flags().GetStringArray(flagName)
		if err != nil {
			return err
		}

		for _, flagValue := range flagValues {
			parts := strings.SplitN(flagValue, "=", 2)
			if _, ok := validOptions[parts[0]]; !ok {
				return fmt.Errorf("invalid output option: %s", parts[0])
			}
			if validOptions[parts[0]] != nil {
				for _, v := range strings.Split(parts[1], ",") {
					if !validOptions[parts[0]][v] {
						return fmt.Errorf("invalid value for output option %s: %s", parts[0], v)
					}
				}
			}
		}
		return nil
	}
}

func outputFlagsForCommand(cmd *cobra.Command) outputOpts {
	opts, _ := cmd.Flags().GetStringArray(flagName)
	return parseOutputFlags(opts)
}

type outputOpts map[string][]string

// Set sets the key to value. It replaces any existing
// values.
func (o outputOpts) Set(key, value string) {
	o[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (o outputOpts) Add(key, value string) {
	o[key] = append(o[key], value)
}

func (o outputOpts) IsSet(key string) bool {
	if values, ok := o[key]; ok && len(values) > 0 {
		return true
	}
	return false
}

func parseOutputFlags(in []string) outputOpts {
	o := outputOpts{}
	for _, param := range in {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) == 2 {
			o[parts[0]] = strings.Split(parts[1], ",")
		} else {
			o[parts[0]] = []string{""}
		}
	}
	return o
}

// newTableOutput creates a new tableOutput.
func newTableOutput() *tableOutput {
	return &tableOutput{
		w:             tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0),
		columns:       map[string]bool{},
		fieldMapping:  map[string]fieldOutputFn{},
		fieldAlias:    map[string]string{},
		allowedFields: map[string]bool{},
	}
}

type fieldOutputFn func(obj interface{}) string

type writerFlusher interface {
	io.Writer
	Flush() error
}

// tableOutput is a generic way to format object as a table.
type tableOutput struct {
	w             writerFlusher
	columns       map[string]bool
	fieldMapping  map[string]fieldOutputFn
	fieldAlias    map[string]string
	allowedFields map[string]bool
}

// Columns returns a list of known output columns.
func (o *tableOutput) Columns() (cols []string) {
	for c := range o.columns {
		cols = append(cols, c)
	}
	sort.Strings(cols)
	return
}

// AddFieldAlias overrides the field name to allow custom column headers.
func (o *tableOutput) AddFieldAlias(field, alias string) *tableOutput {
	o.fieldAlias[field] = alias
	return o
}

// AddFieldOutputFn adds a function which handles the output of the specified field.
func (o *tableOutput) AddFieldOutputFn(field string, fn fieldOutputFn) *tableOutput {
	o.fieldMapping[field] = fn
	o.allowedFields[field] = true
	o.columns[field] = true
	return o
}

// AddAllowedFields reads all first level fieldnames of the struct and allows them to be used.
func (o *tableOutput) AddAllowedFields(obj interface{}) *tableOutput {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		panic("AddAllowedFields input must be a struct")
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		k := t.Field(i).Type.Kind()
		if k != reflect.Bool &&
			k != reflect.Float32 &&
			k != reflect.Float64 &&
			k != reflect.String &&
			k != reflect.Int {
			// only allow simple values
			// complex values need to be mapped via a fieldOutputFn
			continue
		}
		o.allowedFields[strings.ToLower(t.Field(i).Name)] = true
		o.allowedFields[fieldName(t.Field(i).Name)] = true
		o.columns[fieldName(t.Field(i).Name)] = true
	}
	return o
}

// RemoveAllowedField removes fields from the allowed list.
func (o *tableOutput) RemoveAllowedField(fields ...string) *tableOutput {
	for _, field := range fields {
		delete(o.allowedFields, field)
		delete(o.columns, field)
	}
	return o
}

// ValidateColumns returns an error if invalid columns are specified.
func (o *tableOutput) ValidateColumns(cols []string) error {
	var invalidCols []string
	for _, col := range cols {
		if _, ok := o.allowedFields[strings.ToLower(col)]; !ok {
			invalidCols = append(invalidCols, col)
		}
	}
	if len(invalidCols) > 0 {
		return fmt.Errorf("invalid table columns: %s", strings.Join(invalidCols, ","))
	}
	return nil
}

// WriteHeader writes the table header.
func (o *tableOutput) WriteHeader(columns []string) {
	var header []string
	for _, col := range columns {
		if alias, ok := o.fieldAlias[col]; ok {
			col = alias
		}
		header = append(header, strings.Replace(strings.ToUpper(col), "_", " ", -1))
	}
	fmt.Fprintln(o.w, strings.Join(header, "\t"))
}

func (o *tableOutput) Flush() error {
	return o.w.Flush()
}

// Write writes a table line.
func (o *tableOutput) Write(columns []string, obj interface{}) {
	data := structs.Map(obj)
	dataL := map[string]interface{}{}
	for key, value := range data {
		dataL[strings.ToLower(key)] = value
	}

	var out []string
	for _, col := range columns {
		colName := strings.ToLower(col)
		if alias, ok := o.fieldAlias[colName]; ok {
			if fn, ok := o.fieldMapping[alias]; ok {
				out = append(out, fn(obj))
				continue
			}
		}
		if fn, ok := o.fieldMapping[colName]; ok {
			out = append(out, fn(obj))
			continue
		}
		if value, ok := dataL[strings.Replace(colName, "_", "", -1)]; ok {
			if value == nil {
				out = append(out, na(""))
				continue
			}
			if b, ok := value.(bool); ok {
				out = append(out, yesno(b))
				continue
			}
			if s, ok := value.(string); ok {
				out = append(out, na(s))
				continue
			}
			out = append(out, fmt.Sprintf("%v", value))
		}
	}
	fmt.Fprintln(o.w, strings.Join(out, "\t"))
}

func fieldName(name string) string {
	r := []rune(name)
	var out []rune
	for i := range r {
		if i > 0 && (unicode.IsUpper(r[i])) && (i+1 < len(r) && unicode.IsLower(r[i+1]) || unicode.IsLower(r[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(r[i]))
	}
	return string(out)
}
