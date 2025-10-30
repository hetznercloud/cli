package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type testDeprecatable struct {
	isDeprecated         bool
	unavailableAfter     time.Time
	deprecationAnnounced time.Time
}

func (t testDeprecatable) IsDeprecated() bool {
	return t.isDeprecated
}

func (t testDeprecatable) UnavailableAfter() time.Time {
	return t.unavailableAfter
}

func (t testDeprecatable) DeprecationAnnounced() time.Time {
	return t.deprecationAnnounced
}

func TestDescribeDeprecation(t *testing.T) {
	time.Local = time.UTC

	type testCase struct {
		deprecatable hcloud.Deprecatable
		expected     string
	}

	dep := testDeprecatable{
		isDeprecated:         true,
		unavailableAfter:     time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
		deprecationAnnounced: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	tests := map[string]testCase{
		"not deprecated": {
			deprecatable: testDeprecatable{
				isDeprecated: false,
			},
			expected: "",
		},
		"deprecated": {
			deprecatable: dep,
			expected: fmt.Sprintf(
				"Deprecation:\n  Announced:\t\t2021-01-01 00:00:00 UTC (%s)\n  Unavailable After:\t2021-12-31 23:59:59 UTC (%s)\n",
				humanize.Time(dep.DeprecationAnnounced()), humanize.Time(dep.UnavailableAfter()),
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			info := util.DescribeDeprecation(test.deprecatable)
			assert.Equal(t, test.expected, info)
		})
	}
}
