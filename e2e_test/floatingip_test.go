//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/assertjson"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestFloatingIP(t *testing.T) {
	t.Parallel()

	floatingIPName := withSuffix("test-floating-ip")

	_, err := createFloatingIP(t, floatingIPName, "")
	require.EqualError(t, err, "type is required")

	_, err = createFloatingIP(t, floatingIPName, "ipv4")
	require.EqualError(t, err, "one of --home-location or --server is required")

	_, err = createFloatingIP(t, floatingIPName, "ipv4", "--server", "non-existing-server")
	require.EqualError(t, err, "server not found: non-existing-server")

	floatingIPId, err := createFloatingIP(t, floatingIPName, "ipv4", "--home-location", TestLocationName)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("labels", func(t *testing.T) {
		t.Run("add-label-non-existing-floating-ip", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "add-label", "non-existing-floating-ip", "foo=bar")
			require.EqualError(t, err, "floating IP not found: non-existing-floating-ip")
			assert.Empty(t, out)
		})

		t.Run("add-label", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "add-label", strconv.FormatInt(floatingIPId, 10), "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to Floating IP %d\n", floatingIPId), out)
		})
	})

	floatingIPName = withSuffix("new-test-floating-ip")

	out, err := runCommand(t, "floating-ip", "update", strconv.FormatInt(floatingIPId, 10), "--name", floatingIPName, "--description", "Some description")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d updated\n", floatingIPId), out)

	out, err = runCommand(t, "floating-ip", "set-rdns", strconv.FormatInt(floatingIPId, 10), "--hostname", "s1.example.com")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)

	out, err = runCommand(t, "floating-ip", "unassign", "non-existing-floating-ip")
	require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
	assert.Empty(t, out)

	out, err = runCommand(t, "floating-ip", "assign", "non-existing-floating-ip", "non-existing-server")
	require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
	assert.Empty(t, out)

	out, err = runCommand(t, "floating-ip", "assign", strconv.FormatInt(floatingIPId, 10), "non-existing-server")
	require.EqualError(t, err, "server not found: non-existing-server")
	assert.Empty(t, out)

	t.Run("enable-protection", func(t *testing.T) {
		t.Run("unknown-protection-level", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "enable-protection", strconv.FormatInt(floatingIPId, 10), "unknown-protection-level")
			require.EqualError(t, err, "unknown protection level: unknown-protection-level")
			assert.Empty(t, out)
		})

		t.Run("non-existing-floating-ip", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "enable-protection", "non-existing-floating-ip", "delete")
			require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
			assert.Empty(t, out)
		})

		t.Run("enable-delete-protection", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "enable-protection", strconv.FormatInt(floatingIPId, 10), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection enabled for floating IP %d\n", floatingIPId), out)
		})
	})

	ipStr, err := runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10), "--output", "format={{.IP}}")
	require.NoError(t, err)
	ipStr = strings.TrimSpace(ipStr)
	assert.Regexp(t, `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`, ipStr)

	out, err = runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10))
	require.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv4
Name:\s+new-test-floating-ip-[0-9a-f]{8}
Description:\s+Some description
Created:.*?
IP:\s+(?:[0-9]{1,3}\.){3}[0-9]{1,3}
Blocked:\s+no
Home Location:\s+[a-z]{3}[0-9]*
Server:
\s+Not assigned
DNS:
\s+(?:[0-9]{1,3}\.){3}[0-9]{1,3}: s1\.example\.com
Protection:
\s+Delete:\s+yes
Labels:
\s+foo: bar
`, out)

	out, err = runCommand(t, "floating-ip", "list", "--output", "columns=id,name,type,ip,dns,server,home,blocked,protection,labels,created,age")
	require.NoError(t, err)
	assert.Regexp(t, `^ID +NAME +TYPE +IP +DNS +SERVER +HOME +BLOCKED +PROTECTION +LABELS +CREATED +AGE
[0-9]+ +new-test-floating-ip-[0-9a-f]{8} +ipv4 +(?:[0-9]{1,3}\.){3}[0-9]{1,3} +s1\.example\.com +- +[a-z]{3}[0-9]* +no +delete +foo=bar.*?
$`, out)

	out, err = runCommand(t, "floating-ip", "list", "-o=json")
	require.NoError(t, err)
	assertjson.Equal(t, []byte(fmt.Sprintf(`
[
  {
    "id": %d,
    "description": "Some description",
    "created": "<ignore-diff>",
    "ip": "%s",
    "type": "ipv4",
    "server": null,
    "dns_ptr": [
      {
        "ip": "%s",
        "dns_ptr": "s1.example.com"
      }
    ],
    "home_location": {
      "id": "<ignore-diff>",
      "name": "%s",
      "description": "<ignore-diff>",
      "country": "DE",
      "city": "<ignore-diff>",
      "latitude": "<ignore-diff>",
      "longitude": "<ignore-diff>",
      "network_zone": "<ignore-diff>"
    },
    "blocked": false,
    "protection": {
      "delete": true
    },
    "labels": {
      "foo": "bar"
    },
    "name": "%s"
  }
]
`, floatingIPId, ipStr, ipStr, TestLocationName, floatingIPName)), []byte(out))

	out, err = runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIPId, 10))
	assert.Regexp(t, `^Floating IP deletion is protected \(protected, [0-9a-f]+\)$`, err.Error())
	assert.Empty(t, out)

	t.Run("disable-protection", func(t *testing.T) {
		t.Run("non-existing-floating-ip", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "disable-protection", "non-existing-floating-ip", "delete")
			require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
			assert.Empty(t, out)
		})

		t.Run("unknown-protection-level", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "disable-protection", strconv.FormatInt(floatingIPId, 10), "unknown-protection-level")
			require.EqualError(t, err, "unknown protection level: unknown-protection-level")
			assert.Empty(t, out)
		})

		t.Run("disable-delete-protection", func(t *testing.T) {
			out, err = runCommand(t, "floating-ip", "disable-protection", strconv.FormatInt(floatingIPId, 10), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection disabled for floating IP %d\n", floatingIPId), out)
		})
	})

	out, err = runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIPId, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIPId), out)

	floatingIPName = withSuffix("test-floating-ipv6")
	floatingIPId, err = createFloatingIP(t, floatingIPName, "ipv6", "--home-location", TestLocationName)
	if err != nil {
		t.Fatal(err)
	}

	out, err = runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10))
	require.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv6
Name:\s+test-floating-ipv6-[0-9a-f]{8}
Description:\s+-
Created:.*?
IP:\s+[0-9a-f]+:[0-9a-f]+:[0-9a-f]+:[0-9a-f]+::\/64
Blocked:\s+no
Home Location:\s+[a-z]{3}[0-9]*
Server:
\s+Not assigned
DNS:
\s+No reverse DNS entries
Protection:
\s+Delete:\s+no
Labels:
\s+No labels
`, out)

	out, err = runCommand(t, "floating-ip", "describe", strconv.FormatInt(floatingIPId, 10), "--output", "format={{.IP}}")
	require.NoError(t, err)
	out = strings.TrimSpace(out)
	ipv6 := net.ParseIP(out)
	if ipv6 != nil {
		out, err = runCommand(t, "floating-ip", "set-rdns", strconv.FormatInt(floatingIPId, 10), "--ip", ipv6.String()+"1", "--hostname", "s1.example.com")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)

		out, err = runCommand(t, "floating-ip", "set-rdns", strconv.FormatInt(floatingIPId, 10), "--ip", ipv6.String()+"2", "--hostname", "s2.example.com")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)
	} else {
		t.Errorf("invalid IPv6 address: %s", out)
	}

	out, err = runCommand(t, "floating-ip", "list", "-o", "columns=ip,dns")
	require.NoError(t, err)
	assert.Regexp(t, fmt.Sprintf(`^IP +DNS
%s\/64 +2 entries
`, ipv6), out)

	out, err = runCommand(t, "floating-ip", "delete", strconv.FormatInt(floatingIPId, 10))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIPId), out)
}

func createFloatingIP(t *testing.T, name, ipType string, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.FloatingIP.Delete(context.Background(), &hcloud.FloatingIP{Name: name})
	})

	out, err := runCommand(t, append([]string{"floating-ip", "create", "--name", name, "--type", ipType}, args...)...)
	if err != nil {
		return 0, err
	}

	firstLine := strings.Split(out, "\n")[0]
	if !assert.Regexp(t, `^Floating IP [0-9]+ created$`, firstLine) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[12:len(firstLine)-8], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.FloatingIP.Delete(context.Background(), &hcloud.FloatingIP{ID: id})
	})
	return id, nil
}
