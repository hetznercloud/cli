//go:build e2e

package e2e_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createServer(t *testing.T, name, serverType, image string, args ...string) int {
	out, err := runCommand(t, append([]string{"server", "create", "--name", name, "--type", serverType, "--image", image}, args...)...)
	if err != nil {
		t.Fatal(err)
	}

	firstLine := strings.Split(out, "\n")[0]
	if !assert.Regexp(t, `^Server [0-9]+ created$`, firstLine) {
		t.Fatalf("invalid response: %s", out)
	}

	serverID, err := strconv.Atoi(out[7 : len(firstLine)-8])
	if err != nil {
		t.Fatal(err)
	}
	return serverID
}
