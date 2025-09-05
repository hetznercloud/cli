package storagebox

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var EnableSnapshotPlanCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "enable-snapshot-plan [options] <storage-box>",
			Short: "Enable automatic snapshots for a Storage Box",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.StorageBox().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Int("max-snapshots", 0, "Maximum amount of Snapshots that should be created by this Snapshot Plan")
		_ = cmd.MarkFlagRequired("max-snapshots")

		cmd.Flags().Int("minute", 0, "Minute the Snapshot Plan should be executed on (UTC). Not specified means every minute")
		cmd.Flags().Int("hour", 0, "Hour the Snapshot Plan should be executed on (UTC). Not specified means every hour")
		cmd.Flags().String("day-of-week", "", "Day of the week the Snapshot Plan should be executed on. Not specified means every day")
		cmd.Flags().Int("day-of-month", 0, "Day of the month the Snapshot Plan should be executed on. Not specified means every day")

		_ = cmd.RegisterFlagCompletionFunc("day-of-week", cmpl.SuggestCandidates("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"))

		// TODO: should we add some validation here to avoid footguns? (e.g. backing up every minute)

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		maxSnapshots, _ := cmd.Flags().GetInt("max-snapshots")
		minute, _ := cmd.Flags().GetInt("minute")
		hour, _ := cmd.Flags().GetInt("hour")
		dayOfWeek, _ := cmd.Flags().GetString("day-of-week")
		dayOfMonth, _ := cmd.Flags().GetInt("day-of-month")

		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return err
		}
		if storageBox == nil {
			return fmt.Errorf("Storage Box not found: %s", idOrName)
		}

		opts := hcloud.StorageBoxEnableSnapshotPlanOpts{
			MaxSnapshots: maxSnapshots,
		}

		if cmd.Flags().Changed("minute") {
			opts.Minute = &minute
		}
		if cmd.Flags().Changed("hour") {
			opts.Hour = &hour
		}
		if cmd.Flags().Changed("day-of-week") {
			weekday, err := util.WeekdayFromString(dayOfWeek)
			if err != nil {
				return err
			}
			if weekday == 0 {
				// The API expects 7 for Sunday
				weekday = 7
			}
			opts.DayOfWeek = hcloud.Ptr(int(weekday))
		}
		if cmd.Flags().Changed("day-of-month") {
			opts.DayOfMonth = &dayOfMonth
		}

		action, _, err := s.Client().StorageBox().EnableSnapshotPlan(s, storageBox, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Snapshot Plan enabled for Storage Box %d\n", storageBox.ID)
		return nil
	},
}
