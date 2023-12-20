package testutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// CaptureOutStreams redirects stdout & stderr while running fn and returns the outputs as a string.
// If there's an error during capture, it returns the error, otherwise it returns the error
// returned by fn.
func CaptureOutStreams(fn func() error) (string, string, error) {
	outR, outW, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("capture stdout: %v", err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		return "", "", fmt.Errorf("capture stderr: %v", err)
	}

	origOut, origErr := os.Stdout, os.Stderr
	defer func() {
		os.Stdout = origOut
		os.Stderr = origErr
	}()
	os.Stdout, os.Stderr = outW, errW

	outBuf, errBuf := &bytes.Buffer{}, &bytes.Buffer{}

	var copyErr = make(chan error)
	go func() {
		defer close(copyErr)

		copied, err := io.Copy(outBuf, outR)
		if err != nil {
			copyErr <- fmt.Errorf("capture stdout: %v, copied: %d", err, copied)
			return
		}

		copied, err = io.Copy(errBuf, errR)
		if err != nil {
			copyErr <- fmt.Errorf("capture stderr: %v, copied: %d", err, copied)
			return
		}
	}()

	err = fn()

	if err := outW.Close(); err != nil {
		return "", "", fmt.Errorf("capture stdout close pipe writer: %v", err)
	}

	if err := errW.Close(); err != nil {
		return "", "", fmt.Errorf("capture stderr close pipe writer: %v", err)
	}

	if err := <-copyErr; err != nil {
		return "", "", err
	}

	return outBuf.String(), errBuf.String(), err
}
