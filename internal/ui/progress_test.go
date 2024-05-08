package ui

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgressGroup(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	progressGroup := NewProgressGroup(buffer)
	progress1 := progressGroup.Add("progress 1", "")
	progress2 := progressGroup.Add("progress 2", "")
	progress3 := progressGroup.Add("progress 3", "")

	if err := progressGroup.Start(); err != nil {
		t.Fatal(err)
	}

	progress1.SetSuccess()
	progress3.SetError()
	progress2.SetSuccess()

	if err := progressGroup.Stop(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t,
		`progress 1 ...
progress 2 ...
progress 3 ...
progress 1 ... done
progress 3 ... failed
progress 2 ... done
`,
		buffer.String())
}

func TestProgress(t *testing.T) {
	buffer := bytes.NewBuffer(nil)

	progress1 := NewProgress(buffer, "progress 1", "")
	progress1.Start()
	progress1.SetSuccess()

	assert.Equal(t,
		`progress 1 ...
progress 1 ... done
`,
		buffer.String())
}
