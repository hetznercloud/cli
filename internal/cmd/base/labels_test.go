package base

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/hetznercloud/cli/internal/state"
)

// MockResource represents a mock resource for testing
type MockResource struct {
	ID     int64
	Name   string
	Labels map[string]string
}

// MockLabelCmds creates a mock LabelCmds for testing
func createMockLabelCmds() *LabelCmds[*MockResource] {
	return &LabelCmds[*MockResource]{
		ResourceNameSingular: "Mock Resource",
		ShortDescriptionAdd:  "Add a label to a Mock Resource",
		Fetch: func(s state.State, idOrName string) (*MockResource, error) {
			// Mock implementation - return a resource if name starts with "mock"
			if strings.HasPrefix(idOrName, "mock") {
				return &MockResource{
					ID:     1,
					Name:   idOrName,
					Labels: make(map[string]string),
				}, nil
			}
			return nil, &hcloud.Error{Code: hcloud.ErrorCodeNotFound}
		},
		SetLabels: func(s state.State, resource *MockResource, labels map[string]string) error {
			resource.Labels = labels
			return nil
		},
		GetLabels: func(resource *MockResource) map[string]string {
			return resource.Labels
		},
		GetIDOrName: func(resource *MockResource) string {
			return resource.Name
		},
		FetchBatch: func(s state.State, idOrNames []string) ([]*MockResource, []error) {
			resources := make([]*MockResource, len(idOrNames))
			errors := make([]error, len(idOrNames))

			for i, idOrName := range idOrNames {
				if strings.HasPrefix(idOrName, "mock") {
					resources[i] = &MockResource{
						ID:     int64(i + 1),
						Name:   idOrName,
						Labels: make(map[string]string),
					}
				} else {
					errors[i] = &hcloud.Error{Code: hcloud.ErrorCodeNotFound}
				}
			}

			return resources, errors
		},
		SetLabelsBatch: func(s state.State, resources []*MockResource, labels map[string]string) []error {
			errors := make([]error, len(resources))
			for _, resource := range resources {
				if resource != nil {
					resource.Labels = labels
				}
			}
			return errors
		},
	}
}

func TestLabelCmds_RunAddBatch_Success(t *testing.T) {
	cmds := createMockLabelCmds()

	// Mock state (we won't use it in this test)
	var mockState state.State

	// Mock command
	cmd := &cobra.Command{}

	// Test successful batch operation
	args := []string{"mock1", "mock2", "env=prod", "team=backend"}
	err := cmds.RunAddBatch(mockState, cmd, args)

	// Should not return an error for successful operations
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestLabelCmds_RunAddBatch_PartialFailure(t *testing.T) {
	cmds := createMockLabelCmds()

	// Mock state
	var mockState state.State

	// Mock command
	cmd := &cobra.Command{}

	// Test partial failure (some resources exist, some don't)
	args := []string{"mock1", "nonexistent", "mock2", "env=prod"}
	err := cmds.RunAddBatch(mockState, cmd, args)

	// Should return an error due to partial failure
	if err == nil {
		t.Error("Expected error due to partial failure, got nil")
	}

	// Check that error contains information about failures
	if !strings.Contains(err.Error(), "not_found") {
		t.Errorf("Expected error to contain 'not_found', got: %v", err)
	}
}

func TestLabelCmds_RunAddBatch_NoResources(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	cmd := &cobra.Command{}

	// Test with no resources
	args := []string{"env=prod"}
	err := cmds.RunAddBatch(mockState, cmd, args)

	if err == nil {
		t.Error("Expected error for no resources, got nil")
	}

	if !strings.Contains(err.Error(), "must specify at least one") {
		t.Errorf("Expected error to contain 'must specify at least one', got: %v", err)
	}
}

func TestLabelCmds_RunAddBatch_NoLabels(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	cmd := &cobra.Command{}

	// Test with no labels
	args := []string{"mock1", "mock2"}
	err := cmds.RunAddBatch(mockState, cmd, args)

	if err == nil {
		t.Error("Expected error for no labels, got nil")
	}

	if !strings.Contains(err.Error(), "must specify at least one label") {
		t.Errorf("Expected error to contain 'must specify at least one label', got: %v", err)
	}
}

func TestLabelCmds_RunAddBatch_InvalidLabelFormat(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	cmd := &cobra.Command{}

	// Test with invalid label format
	args := []string{"mock1", "invalid-label"}
	err := cmds.RunAddBatch(mockState, cmd, args)

	if err == nil {
		t.Error("Expected error for invalid label format, got nil")
	}

	if !strings.Contains(err.Error(), "must specify at least one label") {
		t.Errorf("Expected error to contain 'must specify at least one label', got: %v", err)
	}
}

func TestLabelCmds_validateAddLabelBatch_Success(t *testing.T) {
	cmds := createMockLabelCmds()

	cmd := &cobra.Command{}
	args := []string{"mock1", "mock2", "env=prod", "team=backend"}

	err := cmds.validateAddLabelBatch(cmd, args)

	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

func TestLabelCmds_validateAddLabelBatch_NoResources(t *testing.T) {
	cmds := createMockLabelCmds()

	cmd := &cobra.Command{}
	args := []string{"env=prod", "team=backend"}

	err := cmds.validateAddLabelBatch(cmd, args)

	if err == nil {
		t.Error("Expected validation error for no resources, got nil")
	}

	if !strings.Contains(err.Error(), "must specify at least one mock resource") {
		t.Errorf("Expected error to contain 'must specify at least one mock resource', got: %v", err)
	}
}

func TestLabelCmds_validateAddLabelBatch_NoLabels(t *testing.T) {
	cmds := createMockLabelCmds()

	cmd := &cobra.Command{}
	args := []string{"mock1", "mock2"}

	err := cmds.validateAddLabelBatch(cmd, args)

	if err == nil {
		t.Error("Expected validation error for no labels, got nil")
	}

	if !strings.Contains(err.Error(), "must specify at least one label") {
		t.Errorf("Expected error to contain 'must specify at least one label', got: %v", err)
	}
}

func TestLabelCmds_validateAddLabelBatch_InvalidLabel(t *testing.T) {
	cmds := createMockLabelCmds()

	cmd := &cobra.Command{}
	args := []string{"mock1", "invalid-label-format"}

	err := cmds.validateAddLabelBatch(cmd, args)

	if err == nil {
		t.Error("Expected validation error for invalid label, got nil")
	}

	if !strings.Contains(err.Error(), "must specify at least one label") {
		t.Errorf("Expected error to contain 'must specify at least one label', got: %v", err)
	}
}

func TestLabelCmds_validateAddLabelBatch_MixedOrder(t *testing.T) {
	cmds := createMockLabelCmds()

	cmd := &cobra.Command{}
	// Test mixed order of resources and labels
	args := []string{"env=prod", "mock1", "team=backend", "mock2"}

	err := cmds.validateAddLabelBatch(cmd, args)

	if err != nil {
		t.Errorf("Expected no validation error for mixed order, got: %v", err)
	}
}

func TestLabelCmds_processAddBatch_Success(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	idOrNames := []string{"mock1", "mock2"}
	labelsToAdd := map[string]string{
		"env":  "prod",
		"team": "backend",
	}

	results := cmds.processAddBatch(mockState, idOrNames, labelsToAdd, false)

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	for i, result := range results {
		if !result.Success {
			t.Errorf("Expected result %d to be successful, got error: %v", i, result.Error)
		}
		if result.IDOrName != idOrNames[i] {
			t.Errorf("Expected IDOrName %s, got %s", idOrNames[i], result.IDOrName)
		}
	}
}

func TestLabelCmds_processAddBatch_ResourceNotFound(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	idOrNames := []string{"mock1", "nonexistent"}
	labelsToAdd := map[string]string{"env": "prod"}

	results := cmds.processAddBatch(mockState, idOrNames, labelsToAdd, false)

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// First result should be successful
	if !results[0].Success {
		t.Errorf("Expected first result to be successful, got error: %v", results[0].Error)
	}

	// Second result should fail
	if results[1].Success {
		t.Error("Expected second result to fail, but it was successful")
	}

	if results[1].Error == nil {
		t.Error("Expected second result to have an error, but it was nil")
	}
}

func TestLabelCmds_displayAddResults(t *testing.T) {
	cmds := createMockLabelCmds()

	cmd := &cobra.Command{}
	results := []BatchLabelResult{
		{IDOrName: "mock1", Success: true, LabelsAdded: []string{"env", "team"}},
		{IDOrName: "mock2", Success: true, LabelsAdded: []string{"env", "team"}},
	}
	labelKeys := []string{"env", "team"}

	// Capture output by redirecting cmd's output
	// This is a simplified test - in a real scenario you'd capture the output
	cmds.displayAddResults(cmd, results, labelKeys)

	// Test with failures
	resultsWithFailures := []BatchLabelResult{
		{IDOrName: "mock1", Success: true, LabelsAdded: []string{"env"}},
		{IDOrName: "nonexistent", Success: false, Error: &hcloud.Error{Code: hcloud.ErrorCodeNotFound}},
	}

	cmds.displayAddResults(cmd, resultsWithFailures, []string{"env"})
}

func TestLabelCmds_fetchResourcesBatch_Success(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	idOrNames := []string{"mock1", "mock2"}

	resources, errors := cmds.fetchResourcesBatch(mockState, idOrNames)

	if len(resources) != 2 {
		t.Errorf("Expected 2 resources, got %d", len(resources))
	}

	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	for i, resource := range resources {
		if resource == nil {
			t.Errorf("Expected resource at index %d to not be nil", i)
		}
		if errors[i] != nil {
			t.Errorf("Expected no error at index %d, got: %v", i, errors[i])
		}
	}
}

func TestLabelCmds_fetchResourcesBatch_WithErrors(t *testing.T) {
	cmds := createMockLabelCmds()

	var mockState state.State
	idOrNames := []string{"mock1", "nonexistent", "mock2"}

	resources, errors := cmds.fetchResourcesBatch(mockState, idOrNames)

	if len(resources) != 3 {
		t.Errorf("Expected 3 resources, got %d", len(resources))
	}

	if len(errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errors))
	}

	// Check first resource (should succeed)
	if resources[0] == nil {
		t.Error("Expected first resource to not be nil")
	}
	if errors[0] != nil {
		t.Errorf("Expected no error for first resource, got: %v", errors[0])
	}

	// Check second resource (should fail)
	if resources[1] != nil {
		t.Error("Expected second resource to be nil due to error")
	}
	if errors[1] == nil {
		t.Error("Expected error for second resource, got nil")
	}

	// Check third resource (should succeed)
	if resources[2] == nil {
		t.Error("Expected third resource to not be nil")
	}
	if errors[2] != nil {
		t.Errorf("Expected no error for third resource, got: %v", errors[2])
	}
}

func TestBatchLabelResult_Structure(t *testing.T) {
	result := BatchLabelResult{
		IDOrName:         "test-resource",
		Resource:         &MockResource{ID: 1, Name: "test"},
		Success:          true,
		Error:            nil,
		LabelsAdded:      []string{"env", "team"},
		LabelsRemoved:    []string{},
		AllLabelsRemoved: false,
	}

	if result.IDOrName != "test-resource" {
		t.Errorf("Expected IDOrName 'test-resource', got '%s'", result.IDOrName)
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.Error != nil {
		t.Errorf("Expected Error to be nil, got: %v", result.Error)
	}

	if len(result.LabelsAdded) != 2 {
		t.Errorf("Expected 2 labels added, got %d", len(result.LabelsAdded))
	}
}

func TestLabelBatchSize_Constant(t *testing.T) {
	if labelBatchSize != 10 {
		t.Errorf("Expected labelBatchSize to be 10, got %d", labelBatchSize)
	}
}

// Benchmark tests for performance validation
func BenchmarkLabelCmds_RunAddBatch(b *testing.B) {
	cmds := createMockLabelCmds()
	var mockState state.State
	cmd := &cobra.Command{}

	// Create a larger batch for benchmarking
	args := make([]string, 0, 20)
	for i := 0; i < 10; i++ {
		args = append(args, "mock"+string(rune(i+'0')))
	}
	args = append(args, "env=prod", "team=backend")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmds.RunAddBatch(mockState, cmd, args)
	}
}

func BenchmarkLabelCmds_processAddBatch(b *testing.B) {
	cmds := createMockLabelCmds()
	var mockState state.State

	idOrNames := make([]string, 10)
	for i := 0; i < 10; i++ {
		idOrNames[i] = "mock" + string(rune(i+'0'))
	}
	labelsToAdd := map[string]string{"env": "prod", "team": "backend"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmds.processAddBatch(mockState, idOrNames, labelsToAdd, false)
	}
}
