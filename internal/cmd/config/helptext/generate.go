package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/hetznercloud/cli/internal/state/config"
)

//go:generate go run $GOFILE

func main() {
	if err := run(); err != nil {
		if _, writeErr := fmt.Fprintln(os.Stderr, err); writeErr != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func run() error {
	if err := generateTable(
		"preferences",
		config.OptionFlagPreference|config.OptionFlagHidden,
		config.OptionFlagPreference,
		table.Row{"sort.<resource>", "Default sorting for resource", "string list", "sort.<resource>", "HCLOUD_SORT_<RESOURCE>", ""},
	); err != nil {
		return err
	}
	return generateTable("other",
		config.OptionFlagPreference|config.OptionFlagHidden,
		0,
	)
}

func generateTable(outFile string, mask, filter config.OptionFlag, extraRows ...table.Row) error {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:             "Description",
			WidthMax:         20,
			WidthMaxEnforcer: text.WrapSoft,
		},
	})

	t.AppendHeader(table.Row{"Option", "Description", "Type", "Config key", "Environment variable", "Flag"})

	var opts []config.IOption
	for _, opt := range config.DefaultOptions() {
		if opt.GetFlags()&mask != filter {
			continue
		}
		opts = append(opts, opt)
	}

	slices.SortFunc(opts, func(a, b config.IOption) int {
		return strings.Compare(a.GetName(), b.GetName())
	})

	for _, opt := range opts {
		typeName, err := getTypeName(opt)
		if err != nil {
			return err
		}
		t.AppendRow(table.Row{opt.GetName(), opt.GetDescription(), typeName, opt.ConfigKey(), opt.EnvVar(), opt.FlagName()})
		t.AppendSeparator()
	}

	for _, row := range extraRows {
		t.AppendRow(row)
		t.AppendSeparator()
	}

	if err := os.WriteFile(outFile+".txt", []byte(t.Render()+"\n"), 0644); err != nil { //nolint:gosec
		return fmt.Errorf("write text table: %w", err)
	}

	if err := os.WriteFile(outFile+".md", []byte(escapeString(t.RenderMarkdown())+"\n"), 0644); err != nil { //nolint:gosec
		return fmt.Errorf("write Markdown table: %w", err)
	}
	return nil
}

func getTypeName(opt config.IOption) (string, error) {
	switch t := opt.T().(type) {
	case bool:
		return "boolean", nil
	case int:
		return "integer", nil
	case string:
		return "string", nil
	case time.Duration:
		return "duration", nil
	case []string:
		return "string list", nil
	default:
		return "", fmt.Errorf("missing type name for %T", t)
	}
}

func escapeString(s string) string {
	replacer := strings.NewReplacer("<", "\\<", ">", "\\>", "_", "\\_")
	return replacer.Replace(s)
}
