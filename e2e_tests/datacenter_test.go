package e2e_tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatacenter(t *testing.T) {
	t.Parallel()

	out, err := runCommand(t, "datacenter", "list")
	assert.NoError(t, err)
	assert.Regexp(t, `ID +NAME +DESCRIPTION +LOCATION
([0-9]+ +[a-z0-9\-]+ +[a-zA-Z0-9\- ]+ +[a-z0-9\-]+\n)+`, out)

	out, err = runCommand(t, "datacenter", "list", "-o=json")
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal([]byte(out), new([]any)))

	out, err = runCommand(t, "datacenter", "describe", "123456")
	assert.EqualError(t, err, "datacenter not found: 123456")
	assert.Empty(t, out)

	out, err = runCommand(t, "datacenter", "describe", "2")
	assert.NoError(t, err)
	assert.Regexp(t, `ID:\s+[0-9]+
Name:\s+[a-z0-9\-]+
Description:\s+[a-zA-Z0-9\- ]+
Location:
 +Name:\s+[a-z0-9]+
 +Description:\s+[a-zA-Z0-9\- ]+
 +Country:\s+[A-Z]{2}
 +City:\s+[A-Za-z]+
 +Latitude:\s+[0-9\.]+
 +Longitude:\s+[0-9\.]+
Server Types:
 +Available:
(\s+- ID:\s+[0-9]+
\s+Name:\s+[a-z0-9]+
\s+Description:\s[A-Za-z0-9 ]+)+
 +Supported:
(\s+- ID:\s+[0-9]+
\s+Name:\s+[a-z0-9]+
\s+Description:\s[A-Za-z0-9 ]+)+
`, out)
}
