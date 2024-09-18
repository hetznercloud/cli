//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func createServer(t *testing.T, name, serverType, image string, args ...string) (int, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _, _ = client.Server.DeleteWithResult(context.Background(), &hcloud.Server{Name: name})
	})

	out, err := runCommand(t, append([]string{"server", "create", "--name", name, "--type", serverType, "--image", image}, args...)...)
	if err != nil {
		return 0, err
	}

	firstLine := strings.Split(out, "\n")[0]
	if !assert.Regexp(t, `^Server [0-9]+ created$`, firstLine) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.Atoi(out[7 : len(firstLine)-8])
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _, _ = client.Server.DeleteWithResult(context.Background(), &hcloud.Server{ID: int64(id)})
	})
	return id, nil
}
