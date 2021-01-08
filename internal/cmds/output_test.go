package cmds

import (
	"bytes"
	"strings"
	"testing"
)

type writerFlusherStub struct {
	bytes.Buffer
}

func (s writerFlusherStub) Flush() error {
	return nil
}

type testFieldsStruct struct {
	Name   string
	Number int
}

func TestTableOutput(t *testing.T) {
	var wfs writerFlusherStub
	to := newTableOutput()
	to.w = &wfs

	t.Run("AddAllowedFields", func(t *testing.T) {
		to.AddAllowedFields(testFieldsStruct{})
		if _, ok := to.allowedFields["name"]; !ok {
			t.Error("name should be a allowed field")
		}
	})
	t.Run("AddFieldAlias", func(t *testing.T) {
		to.AddFieldAlias("leeroy_jenkins", "leeroy jenkins")
		if alias, ok := to.fieldAlias["leeroy_jenkins"]; !ok || alias != "leeroy jenkins" {
			t.Errorf("leeroy_jenkins alias should be 'leeroy jenkins', is: %v", alias)
		}
	})
	t.Run("AddFieldOutputFn", func(t *testing.T) {
		to.AddFieldOutputFn("leeroy jenkins", fieldOutputFn(func(obj interface{}) string {
			return "LEEROY JENKINS!!!"
		}))
		if _, ok := to.fieldMapping["leeroy jenkins"]; !ok {
			t.Errorf("'leeroy jenkins' field output fn should be set")
		}
	})
	t.Run("ValidateColumns", func(t *testing.T) {
		err := to.ValidateColumns([]string{"non-existent", "NAME"})
		if err == nil ||
			strings.Contains(err.Error(), "name") ||
			!strings.Contains(err.Error(), "non-existent") {
			t.Errorf("error should contain 'non-existent' but not 'name': %v", err)
		}
	})
	t.Run("WriteHeader", func(t *testing.T) {
		to.WriteHeader([]string{"leeroy_jenkins", "name"})
		if wfs.String() != "LEEROY JENKINS\tNAME\n" {
			t.Errorf("written header should be 'LEEROY JENKINS\\tNAME\\n', is: %q", wfs.String())
		}
		wfs.Reset()
	})
	t.Run("WriteLine", func(t *testing.T) {
		to.Write([]string{"leeroy_jenkins", "name", "number"}, &testFieldsStruct{"test123", 1000000000})
		if wfs.String() != "LEEROY JENKINS!!!\ttest123\t1000000000\n" {
			t.Errorf("written line should be 'LEEROY JENKINS!!!\\ttest123\\t1000000000\\n', is: %q", wfs.String())
		}
		wfs.Reset()
	})
	t.Run("Columns", func(t *testing.T) {
		if len(to.Columns()) != 3 {
			t.Errorf("unexpected number of columns: %v", to.Columns())
		}
	})
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
