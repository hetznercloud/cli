package cmpl_test

import (
	"testing"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSuggestCandidates(t *testing.T) {
	tests := []struct {
		name       string
		cs         []string
		toComplete string
		sug        []string
		d          cobra.ShellCompDirective
	}{
		{
			name: "no prefix available",
			cs:   []string{"yaml", "json", "toml"},
			sug:  []string{"yaml", "json", "toml"},
		},
		{
			name:       "prefix available",
			cs:         []string{"a", "aa", "aaa", "bbb"},
			toComplete: "aa",
			sug:        []string{"aa", "aaa"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := cmpl.SuggestCandidates(tt.cs...)
			sug, d := f(nil, nil, tt.toComplete)
			assert.Equal(t, tt.sug, sug)
			assert.Equal(t, tt.d, d)
		})
	}
}

func TestSuggestArgs(t *testing.T) {
	tests := []struct {
		name string
		vfs  []func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)
		args []string
		sug  []string
		d    cobra.ShellCompDirective
	}{
		{
			name: "suggest first argument but no vfs provided",
		},
		{
			name: "suggest second argument but no vfs provided",
			args: []string{"aaaa"},
		},
		{
			name: "suggest the only argument",
			vfs: []func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective){
				cmpl.SuggestCandidates("aaaa"),
			},
			sug: []string{"aaaa"},
		},
		{
			name: "suggest the second of three possible arguments",
			vfs: []func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective){
				cmpl.SuggestCandidates("aaaa"),
				cmpl.SuggestCandidates("bbbb"),
				cmpl.SuggestCandidates("cccc"),
			},
			args: []string{"aaaa"},
			sug:  []string{"bbbb"},
		},
		{
			name: "skip suggestions using SuggestNothing",
			vfs: []func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective){
				cmpl.SuggestCandidates("aaaa"),
				cmpl.SuggestNothing(),
				cmpl.SuggestCandidates("cccc"),
			},
			args: []string{"aaaa"},
		},
		{
			name: "skip suggestions using nil",
			vfs: []func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective){
				cmpl.SuggestCandidates("aaaa"),
				nil,
				cmpl.SuggestCandidates("cccc"),
			},
			args: []string{"aaaa"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := cmpl.SuggestArgs(tt.vfs...)
			sug, d := f(nil, tt.args, "")
			assert.Equal(t, tt.sug, sug)
			assert.Equal(t, tt.d, d)
		})
	}
}
