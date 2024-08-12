package testutil

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state"
)

type TestableCommand interface {
	CobraCommand(state.State) *cobra.Command
}

type TestCase struct {
	Args          []string
	PreRun        func(t *testing.T, fx *Fixture)
	ExpOut        string
	ExpOutType    DataType
	ExpErrOut     string
	ExpErrOutType DataType
	ExpErr        string
}

type DataType string

const (
	DataTypeText DataType = "text"
	DataTypeJSON DataType = "json"
	DataTypeYAML DataType = "yaml"
)

func (dt DataType) test(t *testing.T, expected string, actual string, _ ...any) bool {
	switch dt {
	case DataTypeJSON:
		return assert.JSONEq(t, expected, actual)
	case DataTypeYAML:
		if json.Valid([]byte(actual)) {
			t.Error("expected YAML, but got valid JSON")
			return false
		}
		return assert.YAMLEq(t, expected, actual)
	default:
		return assert.Equal(t, expected, actual)
	}
}

func TestCommand(t *testing.T, cmd TestableCommand, cases map[string]TestCase) {
	for name, testCase := range cases {
		if testCase.ExpOutType == "" {
			testCase.ExpOutType = DataTypeText
		}
		if testCase.ExpErrOutType == "" {
			testCase.ExpErrOutType = DataTypeText
		}

		t.Run(name, func(t *testing.T) {
			fx := NewFixture(t)
			defer fx.Finish()

			rootCmd := cli.NewRootCommand(fx.State())
			fx.ExpectEnsureToken()

			rootCmd.AddCommand(cmd.CobraCommand(fx.State()))

			if testCase.PreRun != nil {
				testCase.PreRun(t, fx)
			}

			out, errOut, err := fx.Run(rootCmd, testCase.Args)

			if testCase.ExpErr != "" {
				assert.EqualError(t, err, testCase.ExpErr)
			} else {
				assert.NoError(t, err)
			}
			testCase.ExpOutType.test(t, testCase.ExpOut, out)
			testCase.ExpErrOutType.test(t, testCase.ExpErrOut, errOut)
		})
	}
}
