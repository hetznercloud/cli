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
	generateTable(
		"preferences",
		config.OptionFlagPreference|config.OptionFlagHidden,
		config.OptionFlagPreference,
		table.Row{"sort.<resource>", "Default sorting for resource", "string list", "sort.<resource>", "HCLOUD_SORT_<RESOURCE>", ""},
	)
	generateTable("other",
		config.OptionFlagPreference|config.OptionFlagHidden,
		0,
	)
}

func generateTable(outFile string, mask, filter config.OptionFlag, extraRows ...table.Row) {
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
	for _, opt := range config.Options {
		if opt.GetFlags()&mask != filter {
			continue
		}
		opts = append(opts, opt)
	}

	slices.SortFunc(opts, func(a, b config.IOption) int {
		return strings.Compare(a.GetName(), b.GetName())
	})

	for _, opt := range opts {
		t.AppendRow(table.Row{opt.GetName(), opt.GetDescription(), getTypeName(opt), opt.ConfigKey(), opt.EnvVar(), opt.FlagName()})
		t.AppendSeparator()
	}

	for _, row := range extraRows {
		t.AppendRow(row)
		t.AppendSeparator()
	}

	err := os.WriteFile(outFile+".txt", []byte(t.Render()+"\n"), 0644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outFile+".md", []byte(t.RenderMarkdown()+"\n"), 0644)
	if err != nil {
		panic(err)
	}
}

func getTypeName(opt config.IOption) string {
	switch t := opt.T().(type) {
	case bool:
		return "boolean"
	case int:
		return "integer"
	case string:
		return "string"
	case time.Duration:
		return "duration"
	case []string:
		return "string list"
	default:
		panic(fmt.Sprintf("missing type name for %T", t))
	}
}
