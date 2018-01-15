package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

var validOutputOptsKeys = map[string]bool{
	"columns":  true,
	"noheader": true,
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

func parseOutputOpts(in []string) (outputOpts, error) {
	o := outputOpts{}
	for _, param := range in {
		parts := strings.SplitN(param, "=", 2)
		var key string
		var values []string
		if len(parts) == 2 {
			key = parts[0]
			values = strings.Split(parts[1], ",")
		} else {
			key = param
			values = []string{""}
		}
		if _, ok := validOutputOptsKeys[key]; !ok {
			return o, fmt.Errorf("invalid output key: %s", key)
		}
		o[key] = values
	}
	return o, nil
}

// newTableOutput creates a new tableOutput.
func newTableOutput() *tableOutput {
	return &tableOutput{
		w:             tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0),
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
	fieldMapping  map[string]fieldOutputFn
	fieldAlias    map[string]string
	allowedFields map[string]bool
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
	return o
}

// AddAllowedFields reads all first level fieldnames of the struct and allowes them to be used.
func (o *tableOutput) AddAllowedFields(obj interface{}) *tableOutput {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		panic("AddAllowedFields input must be a struct")
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Ptr {
			continue
		}
		o.allowedFields[strings.ToLower(t.Field(i).Name)] = true
	}
	return o
}

// RemoveAllowedField removes fields from the allowed list.
func (o *tableOutput) RemoveAllowedField(fields ...string) *tableOutput {
	for _, field := range fields {
		delete(o.allowedFields, field)
	}
	return o
}

// ValidateColumns returns an error if invalid columns are specified.
func (o *tableOutput) ValidateColumns(cols []string) error {
	var invalidCols []string
	for _, col := range cols {
		if _, ok := o.allowedFields[col]; !ok {
			invalidCols = append(invalidCols, col)
		}
	}
	if len(invalidCols) > 0 {
		return fmt.Errorf("invalid table columns: %s", strings.Join(invalidCols, ","))
	}
	return nil
}

// WriteHeader writes the table header.
func (o *tableOutput) WriteHeader(collumns []string) {
	var header []string
	for _, col := range collumns {
		if alias, ok := o.fieldAlias[col]; ok {
			col = alias
		}
		header = append(header, strings.ToUpper(col))
	}
	fmt.Fprintln(o.w, strings.Join(header, "\t"))
}

func (o *tableOutput) Flush() error {
	return o.w.Flush()
}

// Write writes a table line.
func (o *tableOutput) Write(collumns []string, obj interface{}) {
	data := map[string]interface{}{}
	objJSON, _ := json.Marshal(obj)
	json.Unmarshal(objJSON, &data)

	dataL := map[string]interface{}{}
	for key, value := range data {
		dataL[strings.ToLower(key)] = value
	}

	var out []string
	for _, col := range collumns {
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
		if value, ok := dataL[colName]; ok {
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
