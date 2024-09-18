//go:build e2e

package e2e_test

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swaggest/assertjson"
)

func TestFloatingIP(t *testing.T) {
	t.Parallel()

	_, err := createFloatingIP(t, "")
	require.EqualError(t, err, "type is required")

	_, err = createFloatingIP(t, "ipv4")
	require.EqualError(t, err, "one of --home-location or --server is required")

	_, err = createFloatingIP(t, "ipv4", "--server", "non-existing-server")
	require.EqualError(t, err, "server not found: non-existing-server")

	floatingIPId, err := createFloatingIP(t, "ipv4", "--home-location", "fsn1")
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
			out, err := runCommand(t, "floating-ip", "add-label", strconv.Itoa(floatingIPId), "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to Floating IP %d\n", floatingIPId), out)
		})
	})

	out, err := runCommand(t, "floating-ip", "update", strconv.Itoa(floatingIPId), "--name", "new-test-floating-ip", "--description", "Some description")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d updated\n", floatingIPId), out)

	out, err = runCommand(t, "floating-ip", "set-rdns", strconv.Itoa(floatingIPId), "--hostname", "s1.example.com")
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)

	out, err = runCommand(t, "floating-ip", "unassign", "non-existing-floating-ip")
	require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
	assert.Empty(t, out)

	out, err = runCommand(t, "floating-ip", "assign", "non-existing-floating-ip", "non-existing-server")
	require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
	assert.Empty(t, out)

	out, err = runCommand(t, "floating-ip", "assign", strconv.Itoa(floatingIPId), "non-existing-server")
	require.EqualError(t, err, "server not found: non-existing-server")
	assert.Empty(t, out)

	t.Run("enable-protection", func(t *testing.T) {
		t.Run("unknown-protection-level", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "enable-protection", strconv.Itoa(floatingIPId), "unknown-protection-level")
			require.EqualError(t, err, "unknown protection level: unknown-protection-level")
			assert.Empty(t, out)
		})

		t.Run("non-existing-floating-ip", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "enable-protection", "non-existing-floating-ip", "delete")
			require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
			assert.Empty(t, out)
		})

		t.Run("enable-delete-protection", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "enable-protection", strconv.Itoa(floatingIPId), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection enabled for floating IP %d\n", floatingIPId), out)
		})
	})

	ipStr, err := runCommand(t, "floating-ip", "describe", strconv.Itoa(floatingIPId), "--output", "format={{.IP}}")
	require.NoError(t, err)
	ipStr = strings.TrimSpace(ipStr)
	assert.Regexp(t, `^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`, ipStr)

	out, err = runCommand(t, "floating-ip", "describe", strconv.Itoa(floatingIPId))
	require.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv4
Name:\s+new-test-floating-ip
Description:\s+Some description
Created:.*?
IP:\s+(?:[0-9]{1,3}\.){3}[0-9]{1,3}
Blocked:\s+no
Home Location:\s+fsn1
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
[0-9]+ +new-test-floating-ip +ipv4 +(?:[0-9]{1,3}\.){3}[0-9]{1,3} +s1\.example\.com +- +fsn1 +no +delete +foo=bar.*?
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
      "id": 1,
      "name": "fsn1",
      "description": "Falkenstein DC Park 1",
      "country": "DE",
      "city": "Falkenstein",
      "latitude": 50.47612,
      "longitude": 12.370071,
      "network_zone": "eu-central"
    },
    "blocked": false,
    "protection": {
      "delete": true
    },
    "labels": {
      "foo": "bar"
    },
    "name": "new-test-floating-ip"
  }
]
`, floatingIPId, ipStr, ipStr)), []byte(out))

	out, err = runCommand(t, "floating-ip", "delete", strconv.Itoa(floatingIPId))
	assert.Regexp(t, `^Floating IP deletion is protected \(protected, [0-9a-f]+\)$`, err.Error())
	assert.Empty(t, out)

	t.Run("disable-protection", func(t *testing.T) {
		t.Run("non-existing-floating-ip", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "disable-protection", "non-existing-floating-ip", "delete")
			require.EqualError(t, err, "Floating IP not found: non-existing-floating-ip")
			assert.Empty(t, out)
		})

		t.Run("unknown-protection-level", func(t *testing.T) {
			out, err := runCommand(t, "floating-ip", "disable-protection", strconv.Itoa(floatingIPId), "unknown-protection-level")
			require.EqualError(t, err, "unknown protection level: unknown-protection-level")
			assert.Empty(t, out)
		})

		t.Run("disable-delete-protection", func(t *testing.T) {
			out, err = runCommand(t, "floating-ip", "disable-protection", strconv.Itoa(floatingIPId), "delete")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Resource protection disabled for floating IP %d\n", floatingIPId), out)
		})
	})

	out, err = runCommand(t, "floating-ip", "delete", strconv.Itoa(floatingIPId))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIPId), out)

	floatingIPId, err = createFloatingIP(t, "ipv6", "--home-location", "fsn1")
	if err != nil {
		t.Fatal(err)
	}

	out, err = runCommand(t, "floating-ip", "describe", strconv.Itoa(floatingIPId))
	require.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Type:\s+ipv6
Name:\s+test-floating-ip
Description:\s+-
Created:.*?
IP:\s+[0-9a-f]+:[0-9a-f]+:[0-9a-f]+:[0-9a-f]+::\/64
Blocked:\s+no
Home Location:\s+fsn1
Server:
\s+Not assigned
DNS:
\s+No reverse DNS entries
Protection:
\s+Delete:\s+no
Labels:
\s+No labels
`, out)

	out, err = runCommand(t, "floating-ip", "describe", strconv.Itoa(floatingIPId), "--output", "format={{.IP}}")
	require.NoError(t, err)
	out = strings.TrimSpace(out)
	ipv6 := net.ParseIP(out)
	if ipv6 != nil {
		out, err = runCommand(t, "floating-ip", "set-rdns", strconv.Itoa(floatingIPId), "--ip", ipv6.String()+"1", "--hostname", "s1.example.com")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Reverse DNS of Floating IP %d changed\n", floatingIPId), out)

		out, err = runCommand(t, "floating-ip", "set-rdns", strconv.Itoa(floatingIPId), "--ip", ipv6.String()+"2", "--hostname", "s2.example.com")
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

	out, err = runCommand(t, "floating-ip", "delete", strconv.Itoa(floatingIPId))
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Floating IP %d deleted\n", floatingIPId), out)
}

func createFloatingIP(t *testing.T, ipType string, args ...string) (int, error) {
	out, err := runCommand(t, append([]string{"floating-ip", "create", "--name", "test-floating-ip", "--type", ipType}, args...)...)
	if err != nil {
		return 0, err
	}

	firstLine := strings.Split(out, "\n")[0]
	if !assert.Regexp(t, `^Floating IP [0-9]+ created$`, firstLine) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.Atoi(out[12 : len(firstLine)-8])
	if err != nil {
		return 0, err
	}
	return id, nil
}
