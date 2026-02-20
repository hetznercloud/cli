package base

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

// labelBatchSize is the batch size when processing multiple resources in parallel.
const labelBatchSize = 10

// BatchLabelResult represents the result of a batch label operation for a single resource.
type BatchLabelResult struct {
	IDOrName      string
	Resource      interface{}
	Success       bool
	Error         error
	LabelsAdded   []string
	LabelsRemoved []string
	AllLabelsRemoved bool
}

// LabelCmds allows defining commands for adding labels to resources.
type LabelCmds[T any] struct {
	ResourceNameSingular   string
	ShortDescriptionAdd    string
	ShortDescriptionRemove string
	NameSuggestions        func(client hcapi2.Client) func() []string
	LabelKeySuggestions    func(client hcapi2.Client) func(idOrName string) []string
	Fetch                  func(s state.State, idOrName string) (T, error)

	// Can be set in case the resource has more than a single identifier that is used in the positional arguments.
	// See [LabelCmds.PositionalArgumentOverride].
	FetchWithArgs func(s state.State, args []string) (T, error)
	SetLabels     func(s state.State, resource T, labels map[string]string) error
	GetLabels     func(resource T) map[string]string
	GetIDOrName   func(resource T) string

	// In case the resource does not have a single identifier that matches [LabelCmds.ResourceNameSingular], this field
	// can be set to define the list of positional arguments.
	// For example, passing:
	//     []string{"a", "b", "c"}
	// Would result in the usage string:
	//     <a> <b> <c>
	PositionalArgumentOverride []string

	// Can be set if the default [LabelCmds.NameSuggestions] is not enough. This is usually the case when
	// [LabelCmds.FetchWithArgs] and [LabelCmds.PositionalArgumentOverride] is being used.
	//
	// If this is being set [LabelCmds.LabelKeySuggestions] is ignored and its functionality must be
	// provided as part of the [LabelCmds.ValidArgsFunction].
	ValidArgsFunction func(client hcapi2.Client) []cobra.CompletionFunc

	// Experimental is a function that will be used to mark the command as experimental.
	Experimental func(state.State, *cobra.Command) *cobra.Command

	// Batch operation support
	FetchBatch    func(s state.State, idOrNames []string) ([]T, []error)
	SetLabelsBatch func(s state.State, resources []T, labels map[string]string) []error
}

// AddCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds[T]) AddCobraCommand(s state.State) *cobra.Command {
	var suggestArgs []cobra.CompletionFunc
	switch {
	case lc.NameSuggestions != nil:
		suggestArgs = append(suggestArgs,
			cmpl.SuggestCandidatesF(lc.NameSuggestions(s.Client())),
		)
	case lc.ValidArgsFunction != nil:
		suggestArgs = append(suggestArgs, lc.ValidArgsFunction(s.Client())...)
	default:
		log.Fatalf("label command %s is missing ValidArgsFunction or NameSuggestions", lc.ResourceNameSingular)
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("add-label [--overwrite] %s... <label>...", positionalArguments(lc.ResourceNameSingular, lc.PositionalArgumentOverride)),
		Short:                 lc.ShortDescriptionAdd,
		Args:                  util.Validate,
		ValidArgsFunction:     cmpl.SuggestArgs(suggestArgs...),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(lc.validateAddLabelBatch, s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.RunAddBatch(s, cmd, args)
		},
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already (true, false)")

	if lc.Experimental != nil {
		cmd = lc.Experimental(s, cmd)
	}

	return cmd
}

// RunAdd executes an add label command
func (lc *LabelCmds[T]) RunAdd(s state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	var resource T
	var err error

	if lc.FetchWithArgs != nil {
		resource, err = lc.FetchWithArgs(s, args)
	} else {
		resource, err = lc.Fetch(s, args[0])
	}

	if err != nil {
		return err
	}

	labels, idOrName := lc.GetLabels(resource), lc.GetIDOrName(resource)
	toAdd := args[max(1, len(lc.PositionalArgumentOverride)):]

	if labels == nil {
		labels = map[string]string{}
	}

	var keys []string
	for _, label := range toAdd {
		key, val := util.SplitLabelVars(label)
		keys = append(keys, key)

		if _, ok := labels[key]; ok && !overwrite {
			return fmt.Errorf("label %s on %s %s already exists", key, lc.ResourceNameSingular, idOrName)
		}

		labels[key] = val
	}

	if err := lc.SetLabels(s, resource, labels); err != nil {
		return err
	}

	cmd.Printf("Label(s) %s added to %s %s\n", strings.Join(keys, ", "), lc.ResourceNameSingular, idOrName)
	return nil
}

// RunAddBatch executes a batch add label command
func (lc *LabelCmds[T]) RunAddBatch(s state.State, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")

	// Parse arguments to separate resources from labels
	var resourceArgs []string
	var labelArgs []string

	// Separate arguments: those with '=' are labels, others are resources
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			labelArgs = append(labelArgs, arg)
		} else {
			resourceArgs = append(resourceArgs, arg)
		}
	}

	if len(resourceArgs) == 0 {
		return fmt.Errorf("must specify at least one %s", strings.ToLower(lc.ResourceNameSingular))
	}

	if len(labelArgs) == 0 {
		return fmt.Errorf("must specify at least one label")
	}

	// Parse labels to add
	labelsToAdd := make(map[string]string)
	var labelKeys []string
	for _, label := range labelArgs {
		key, val := util.SplitLabelVars(label)
		labelsToAdd[key] = val
		labelKeys = append(labelKeys, key)
	}

	// Process resources in batches
	errs := make([]error, 0, len(resourceArgs))
	results := make([]BatchLabelResult, 0, len(resourceArgs))

	for batch := range slices.Chunk(resourceArgs, labelBatchSize) {
		batchResults := lc.processAddBatch(s, batch, labelsToAdd, overwrite)
		results = append(results, batchResults...)

		// Collect errors
		for _, result := range batchResults {
			if result.Error != nil {
				errs = append(errs, result.Error)
			}
		}
	}

	// Display results
	lc.displayAddResults(cmd, results, labelKeys)

	return errors.Join(errs...)
}

// processAddBatch processes a batch of resources for adding labels
func (lc *LabelCmds[T]) processAddBatch(s state.State, idOrNames []string, labelsToAdd map[string]string, overwrite bool) []BatchLabelResult {
	results := make([]BatchLabelResult, len(idOrNames))

	// Fetch all resources in batch
	resources, fetchErrors := lc.fetchResourcesBatch(s, idOrNames)

	for i, idOrName := range idOrNames {
		results[i] = BatchLabelResult{IDOrName: idOrName}

		if fetchErrors[i] != nil {
			results[i].Error = fetchErrors[i]
			continue
		}

		resource := resources[i]
		if util.IsNil(resource) {
			results[i].Error = fmt.Errorf("%s not found: %s", lc.ResourceNameSingular, idOrName)
			continue
		}

		// Get current labels
		existingLabels := lc.GetLabels(resource)
		if existingLabels == nil {
			existingLabels = make(map[string]string)
		}

		// Check for conflicts
		var conflictKeys []string
		for key := range labelsToAdd {
			if _, exists := existingLabels[key]; exists && !overwrite {
				conflictKeys = append(conflictKeys, key)
			}
		}

		if len(conflictKeys) > 0 {
			results[i].Error = fmt.Errorf("label(s) %s on %s %s already exist",
				strings.Join(conflictKeys, ", "), lc.ResourceNameSingular, idOrName)
			continue
		}

		// Merge new labels
		for key, value := range labelsToAdd {
			existingLabels[key] = value
		}

		// Apply labels
		if err := lc.SetLabels(s, resource, existingLabels); err != nil {
			results[i].Error = err
			continue
		}

		results[i].Success = true
		results[i].Resource = resource
		results[i].LabelsAdded = slices.Collect(maps.Keys(labelsToAdd))
	}

	return results
}

// fetchResourcesBatch fetches multiple resources in parallel
func (lc *LabelCmds[T]) fetchResourcesBatch(s state.State, idOrNames []string) ([]T, []error) {
	if lc.FetchBatch != nil {
		return lc.FetchBatch(s, idOrNames)
	}

	// Fallback to individual fetching with goroutines
	resources := make([]T, len(idOrNames))
	errors := make([]error, len(idOrNames))

	var wg sync.WaitGroup
	for i, idOrName := range idOrNames {
		wg.Add(1)
		go func(idx int, id string) {
			defer wg.Done()
			resource, err := lc.Fetch(s, id)
			resources[idx] = resource
			errors[idx] = err
		}(i, idOrName)
	}
	wg.Wait()

	return resources, errors
}

// displayAddResults displays the results of a batch add operation
func (lc *LabelCmds[T]) displayAddResults(cmd *cobra.Command, results []BatchLabelResult, labelKeys []string) {
	successful := 0
	failed := 0

	for _, result := range results {
		if result.Success {
			successful++
			cmd.Printf("✓ %s: Labels %s added\n",
				result.IDOrName,
				strings.Join(labelKeys, ", "))
		} else {
			failed++
			cmd.Printf("✗ %s: %v\n", result.IDOrName, result.Error)
		}
	}

	// Summary for multiple resources
	if len(results) > 1 {
		cmd.Printf("\nSummary: %d successful, %d failed\n", successful, failed)
	}
}

func (lc *LabelCmds[T]) validateAddLabel(_ *cobra.Command, args []string) error {
	numPosArgs := max(1, len(lc.PositionalArgumentOverride))

	for _, label := range args[numPosArgs:] {
		if len(util.SplitLabel(label)) != 2 {
			return fmt.Errorf("invalid label: %s", label)
		}
	}

	return nil
}

// validateAddLabelBatch validates arguments for batch add label operations
func (lc *LabelCmds[T]) validateAddLabelBatch(_ *cobra.Command, args []string) error {
	var resourceArgs []string
	var labelArgs []string

	// Separate arguments: those with '=' are labels, others are resources
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			labelArgs = append(labelArgs, arg)
		} else {
			resourceArgs = append(resourceArgs, arg)
		}
	}

	if len(resourceArgs) == 0 {
		return fmt.Errorf("must specify at least one %s", strings.ToLower(lc.ResourceNameSingular))
	}

	if len(labelArgs) == 0 {
		return fmt.Errorf("must specify at least one label")
	}

	// Validate label format for all labels
	for _, label := range labelArgs {
		if len(util.SplitLabel(label)) != 2 {
			return fmt.Errorf("invalid label: %s", label)
		}
	}

	return nil
}

// RemoveCobraCommand creates a command that can be registered with cobra.
func (lc *LabelCmds[T]) RemoveCobraCommand(s state.State) *cobra.Command {
	var suggestArgs []cobra.CompletionFunc
	switch {
	case lc.NameSuggestions != nil:
		suggestArgs = append(suggestArgs,
			cmpl.SuggestCandidatesF(lc.NameSuggestions(s.Client())),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return lc.LabelKeySuggestions(s.Client())(args[0])
			}),
		)
	case lc.ValidArgsFunction != nil:
		suggestArgs = append(suggestArgs, lc.ValidArgsFunction(s.Client())...)
	default:
		log.Fatalf("label command %s is missing ValidArgsFunction or NameSuggestions", lc.ResourceNameSingular)
	}

	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("remove-label %s (--all | <label>...)", positionalArguments(lc.ResourceNameSingular, lc.PositionalArgumentOverride)),
		Short:                 lc.ShortDescriptionRemove,
		Args:                  util.ValidateLenient,
		ValidArgsFunction:     cmpl.SuggestArgs(suggestArgs...),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(lc.validateRemoveLabel, s.EnsureToken),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.RunRemove(s, cmd, args)
		},
	}

	cmd.Flags().BoolP("all", "a", false, "Remove all labels")

	if lc.Experimental != nil {
		cmd = lc.Experimental(s, cmd)
	}

	return cmd
}

// RunRemove executes a remove label command
func (lc *LabelCmds[T]) RunRemove(s state.State, cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	var resource T
	var err error

	if lc.FetchWithArgs != nil {
		resource, err = lc.FetchWithArgs(s, args)
	} else {
		resource, err = lc.Fetch(s, args[0])
	}

	if err != nil {
		return err
	}

	labels, idOrName := lc.GetLabels(resource), lc.GetIDOrName(resource)
	toRemove := args[max(1, len(lc.PositionalArgumentOverride)):]

	if all {
		labels = make(map[string]string)
	} else {
		for _, key := range toRemove {
			if _, ok := labels[key]; !ok {
				return fmt.Errorf("label %s on %s %s does not exist", key, lc.ResourceNameSingular, idOrName)
			}
			delete(labels, key)
		}
	}

	if err := lc.SetLabels(s, resource, labels); err != nil {
		return err
	}

	if all {
		cmd.Printf("All labels removed from %s %s\n", lc.ResourceNameSingular, idOrName)
	} else {
		cmd.Printf("Label(s) %s removed from %s %s\n", strings.Join(toRemove, ", "), lc.ResourceNameSingular, idOrName)
	}

	return nil
}

func (lc *LabelCmds[T]) validateRemoveLabel(cmd *cobra.Command, args []string) error {
	all, _ := cmd.Flags().GetBool("all")

	numPosArgs := max(1, len(lc.PositionalArgumentOverride))

	if all && len(args) > numPosArgs {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) <= numPosArgs {
		return errors.New("must specify a label key when not using --all/-a")
	}

	return nil
}

func positionalArguments(resourceNameSingular string, positionalArgumentOverride []string) string {
	args := []string{util.ToKebabCase(resourceNameSingular)}
	if len(positionalArgumentOverride) > 0 {
		args = slices.Clone(positionalArgumentOverride)
	}

	for i, arg := range args {
		args[i] = fmt.Sprintf("<%s>", arg)
	}

	return strings.Join(args, " ")
}
