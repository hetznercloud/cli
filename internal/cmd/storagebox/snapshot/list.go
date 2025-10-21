package snapshot

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = base.ListCmd[*hcloud.StorageBoxSnapshot, schema.StorageBoxSnapshot]{
	ResourceNamePlural: "Storage Box Snapshots",
	JSONKeyGetByName:   "snapshots",

	DefaultColumns: []string{"id", "name", "description", "size", "is_automatic", "age"},

	ValidArgsFunction: func(client hcapi2.Client) cobra.CompletionFunc {
		return cmpl.SuggestCandidatesF(client.StorageBox().Names)
	},

	PositionalArgumentOverride: []string{"storage-box"},
	SortOption:                 config.OptionSortStorageBoxSnapshot,

	AdditionalFlags: func(cmd *cobra.Command) {
		cmd.Flags().Bool("automatic", false, "Only show automatic snapshots (true, false)")
	},

	FetchWithArgs: func(s state.State, flags *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string, args []string) ([]*hcloud.StorageBoxSnapshot, error) {
		storageBoxIDOrName := args[0]
		isAutomatic, _ := flags.GetBool("automatic")

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		opts := hcloud.StorageBoxSnapshotListOpts{LabelSelector: listOpts.LabelSelector}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		if flags.Changed("automatic") {
			opts.IsAutomatic = &isAutomatic
		}
		return s.Client().StorageBox().AllSnapshotsWithOpts(s, storageBox, opts)
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBoxSnapshot{}).
			AddFieldFn("size", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				return humanize.IBytes(snapshot.Stats.Size)
			}).
			AddFieldFn("size_filesystem", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				return humanize.IBytes(snapshot.Stats.SizeFilesystem)
			}).
			AddFieldFn("labels", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				return util.LabelsToString(snapshot.Labels)
			}).
			AddFieldFn("created", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				return util.Datetime(snapshot.Created)
			}).
			AddFieldFn("age", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				return util.Age(snapshot.Created, time.Now())
			})
	},

	Schema:       hcloud.SchemaFromStorageBoxSnapshot,
	Experimental: experimental.StorageBoxes,
}
