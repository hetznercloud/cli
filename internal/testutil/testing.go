package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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

// MockResponse returns a *hcloud.Response with the given value as JSON body.
func MockResponse[V any](v V) (*hcloud.Response, error) {

	responseBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}

	responseRecorder := httptest.NewRecorder()
	responseRecorder.WriteHeader(200)

	_, err = responseRecorder.Write(responseBytes)
	if err != nil {
		return nil, err
	}

	return &hcloud.Response{
		Response: responseRecorder.Result(),
	}, nil
}
