package output

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type writerStub struct {
	bytes.Buffer
}

func reset(to *Table[any], ws *writerStub) {
	to.w.ResetHeaders()
	to.w.ResetRows()
	to.w.ResetFooters()
	ws.Reset()
}

type testFieldsStruct struct {
	Name   string
	Number int
}

func TestTableOutput(t *testing.T) {
	var ws writerStub
	to := NewTable[any](io.Discard)
	to.out = &ws

	t.Run("AddAllowedFields", func(t *testing.T) {
		to.AddAllowedFields(testFieldsStruct{})

		assert.Contains(t, to.allowedFields, "name")
	})

	t.Run("AddFieldAlias", func(t *testing.T) {
		to.AddFieldAlias("leeroy_jenkins", "leeroy jenkins")

		assert.Contains(t, to.fieldAlias, "leeroy_jenkins")
		assert.Equal(t, "leeroy jenkins", to.fieldAlias["leeroy_jenkins"])
	})

	t.Run("AddFieldOutputFn", func(t *testing.T) {
		to.AddFieldFn("leeroy jenkins", func(any) string {
			return "LEEROY JENKINS!!!"
		})

		assert.Contains(t, to.fieldMapping, "leeroy jenkins")
	})

	t.Run("MarkFieldAsDeprecated", func(t *testing.T) {
		to.MarkFieldAsDeprecated("name", "This field is deprecated")

		assert.Contains(t, to.deprecations, "name")
	})

	t.Run("ValidateColumns", func(t *testing.T) {
		warnings, err := to.ValidateColumns([]string{"non-existent", "NAME"})

		require.ErrorContains(t, err, "non-existent")
		assert.NotContains(t, err.Error(), "name")
		assert.Contains(t, warnings, "This field is deprecated")
	})

	t.Run("WriteHeader", func(t *testing.T) {
		to.WriteHeader([]string{"leeroy_jenkins", "name"})
		require.NoError(t, to.Flush())

		assert.Equal(t, "LEEROY JENKINS   NAME\n", ws.String())

		reset(to, &ws)
	})

	t.Run("WriteLine", func(t *testing.T) {
		require.NoError(t, to.Write(t.Context(), []string{"leeroy_jenkins", "name", "number"}, &testFieldsStruct{"test123", 1000000000}))
		require.NoError(t, to.Flush())

		assert.Equal(t, "LEEROY JENKINS!!!   test123   1000000000\n", ws.String())

		reset(to, &ws)
	})

	t.Run("Columns", func(t *testing.T) {
		if len(to.Columns()) != 3 {
			t.Errorf("unexpected number of columns: %v", to.Columns())
		}
	})
}

func TestAddAllowedFieldsRejectsNonStruct(t *testing.T) {
	to := NewTable[any](io.Discard)
	to.AddAllowedFields("not a struct")

	_, err := to.ValidateColumns(nil)
	require.ErrorContains(t, err, "got string")
}

func TestTableWritePropagatesFieldError(t *testing.T) {
	to := NewTable[any](io.Discard).
		AddFieldFnE("remote", func(context.Context, any) (string, error) {
			return "", errors.New("lookup failed")
		})

	err := to.Write(t.Context(), []string{"remote"}, struct{}{})
	require.ErrorContains(t, err, `render column "remote": lookup failed`)
}

func TestFieldName(t *testing.T) {
	type fixture struct {
		in, out string
	}
	tests := []fixture{
		{"test", "test"},
		{"t", "t"},
		{"T", "t"},
		{"Server", "server"},
		{"BoundTo", "bound_to"},
	}

	for _, test := range tests {
		if f := fieldName(test.in); f != test.out {
			t.Errorf("Unexpected output expected %q, is %q", test.out, f)
		}
	}
}
