package snapshot

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
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

	FetchWithArgs: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, _ []string, args []string) ([]*hcloud.StorageBoxSnapshot, error) {
		storageBoxIDOrName := args[0]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		opts := hcloud.StorageBoxSnapshotListOpts{LabelSelector: listOpts.LabelSelector}
		return s.Client().StorageBox().AllSnapshotsWithOpts(s, storageBox, opts)
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
			AddAllowedFields(hcloud.StorageBoxSnapshot{}).
			AddFieldFn("description", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				return util.OptionalString(snapshot.Description, "-")
			}).
			AddFieldFn("size", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				if snapshot.Stats == nil {
					return "-"
				}
				return humanize.IBytes(snapshot.Stats.Size)
			}).
			AddFieldFn("size_filesystem", func(obj any) string {
				snapshot := obj.(*hcloud.StorageBoxSnapshot)
				if snapshot.Stats == nil {
					return "-"
				}
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

	Schema: hcloud.SchemaFromStorageBoxSnapshot,
}
