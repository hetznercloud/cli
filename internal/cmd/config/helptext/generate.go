package main

import (
	"os"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/hetznercloud/cli/internal/state/config"
)

//go:generate go run $GOFILE

func main() {
	generateTable("preferences.txt", config.OptionFlagPreference, true)
	generateTable("other.txt", config.OptionFlagPreference, false)
}

func generateTable(outFile string, filterFlag config.OptionFlag, hasFlag bool) {
	f, err := os.OpenFile(outFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:             "Description",
			WidthMax:         20,
			WidthMaxEnforcer: text.WrapSoft,
		},
	})

	t.SetOutputMirror(f)
	t.AppendHeader(table.Row{"Option", "Description", "Config key", "Environment variable", "Flag"})

	var opts []config.IOption
	for _, opt := range config.Options {
		if opt.HasFlags(filterFlag) != hasFlag {
			continue
		}
		opts = append(opts, opt)
	}

	slices.SortFunc(opts, func(a, b config.IOption) int {
		return strings.Compare(a.GetName(), b.GetName())
	})

	for _, opt := range opts {
		t.AppendRow(table.Row{opt.GetName(), opt.GetDescription(), opt.ConfigKey(), opt.EnvVar(), opt.FlagName()})
		t.AppendSeparator()
	}

	t.Render()
}
