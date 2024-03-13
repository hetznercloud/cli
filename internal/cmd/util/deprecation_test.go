package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
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

	dep := testDeprecatable{
		isDeprecated:         true,
		unavailableAfter:     time.Date(2021, 12, 31, 23, 59, 59, 0, time.UTC),
		deprecationAnnounced: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	info := util.DescribeDeprecation(dep)
	expected := fmt.Sprintf(
		"Deprecation:\n  Announced:\t\tFri Jan  1 00:00:00 UTC 2021 (%s)\n  Unavailable After:\tFri Dec 31 23:59:59 UTC 2021 (%s)\n",
		humanize.Time(dep.DeprecationAnnounced()), humanize.Time(dep.UnavailableAfter()),
	)
	assert.Equal(t, expected, info)
}
