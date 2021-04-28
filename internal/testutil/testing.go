package testutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// CaptureStdout redirects stdout while running fn and returns the output as a string.
// If there's an error during capture, it returns the error, otherwise it returns the error
// returned by fn.
func CaptureStdout(fn func() error) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("capture stdout: %v", err)
	}

	origOut := os.Stdout
	defer func() {
		os.Stdout = origOut
	}()

	buf := &bytes.Buffer{}
	os.Stdout = w

	copyDone := make(chan struct{})
	var copyErr error
	go func() {
		defer close(copyDone)

		copied, err := io.Copy(buf, r)
		if err != nil {
			copyErr = fmt.Errorf("capture stdout: %v, copied: %d\n", err, copied)
			return
		}
	}()

	err = fn()

	if copyErr != nil {
		return "", copyErr
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("capture stdout close pipe reader: %v", err)
	}

	<-copyDone

	return buf.String(), err
}
