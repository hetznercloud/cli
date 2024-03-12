package util_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
)

func TestYesNo(t *testing.T) {
	assert.Equal(t, "yes", util.YesNo(true))
	assert.Equal(t, "no", util.YesNo(false))
}

func TestNA(t *testing.T) {
	assert.Equal(t, "-", util.NA(""))
	assert.Equal(t, "foo", util.NA("foo"))
}

func TestDatetime(t *testing.T) {
	time.Local = time.UTC
	tm := time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC)
	assert.Equal(t, "Thu Nov 17 15:22:12 UTC 2022", util.Datetime(tm))
}

func TestChainRunE(t *testing.T) {
	var calls int
	f1 := func(_ *cobra.Command, args []string) error {
		calls++
		return nil
	}
	f2 := func(_ *cobra.Command, args []string) error {
		calls++
		return errors.New("error")
	}
	f3 := func(_ *cobra.Command, args []string) error {
		calls++
		return nil
	}

	fn := util.ChainRunE(f1, f2, f3)
	err := fn(nil, nil)

	assert.EqualError(t, err, "error")
	assert.Equal(t, 2, calls)
}

func TestOnlyOneSet(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		ss       []string
		expected bool
	}{
		{
			name:     "only arg emtpy",
			expected: false,
		},
		{
			name:     "only arg non-empty",
			s:        "s",
			expected: true,
		},
		{
			name:     "first arg non-empty, rest empty",
			s:        "s",
			ss:       []string{""},
			expected: true,
		},
		{
			name: "at least one other arg non-empty",
			s:    "s",
			ss:   []string{"", "s"},
		},
		{
			name:     "only one arg non-empty",
			ss:       []string{"", "s"},
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := util.ExactlyOneSet(tt.s, tt.ss...)
			if tt.expected != actual {
				t.Errorf("expected %t; got %t", tt.expected, actual)
			}
		})
	}
}

func TestAge(t *testing.T) {
	tests := []struct {
		name     string
		t        time.Time
		now      time.Time
		expected string
	}{
		{
			name:     "exactly now",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC),
			expected: "just now",
		},
		{
			name:     "within a few milliseconds",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 11, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			expected: "just now",
		},
		{
			name:     "10 seconds",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 22, 22, 21, time.UTC),
			expected: "10s",
		},
		{
			name:     "10 minutes",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 17, 15, 32, 12, 21, time.UTC),
			expected: "10m",
		},
		{
			name:     "24 hours",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 18, 15, 22, 12, 21, time.UTC),
			expected: "1d",
		},
		{
			name:     "25 hours",
			t:        time.Date(2022, 11, 17, 15, 22, 12, 21, time.UTC),
			now:      time.Date(2022, 11, 18, 16, 22, 12, 21, time.UTC),
			expected: "1d",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := util.Age(tt.t, tt.now)
			if tt.expected != actual {
				t.Errorf("expected %s; got %s", tt.expected, actual)
			}
		})
	}
}

func TestSplitLabel(t *testing.T) {
	assert.Equal(t, []string{"foo", "bar"}, util.SplitLabel("foo=bar"))
	assert.Equal(t, []string{"foo"}, util.SplitLabel("foo"))
	assert.Equal(t, []string{""}, util.SplitLabel(""))
}

func TestSplitLabelVars(t *testing.T) {
	var a, b string
	a, b = util.SplitLabelVars("foo=bar")
	assert.Equal(t, "foo", a)
	assert.Equal(t, "bar", b)
	a, b = util.SplitLabelVars("foo")
	assert.Equal(t, "foo", a)
	assert.Equal(t, "", b)
	a, b = util.SplitLabelVars("")
	assert.Equal(t, "", a)
	assert.Equal(t, "", b)
}

func TestLabelsToString(t *testing.T) {
	assert.Contains(t, []string{"foo=bar, baz=qux", "baz=qux, foo=bar"}, util.LabelsToString(map[string]string{
		"foo": "bar",
		"baz": "qux",
	}))
	assert.Equal(t, "foo=bar", util.LabelsToString(map[string]string{
		"foo": "bar",
	}))
	assert.Equal(t, "", util.LabelsToString(map[string]string{}))
}

func TestPrefixLines(t *testing.T) {
	assert.Equal(t, "  foo\n  bar", util.PrefixLines("foo\nbar", "  "))
}

func TestDescribeFormat(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return util.DescribeFormat(struct {
			Foo string
			Bar string
		}{
			Foo: "foo",
			Bar: "bar",
		}, "Foo is: {{.Foo}} Bar is: {{.Bar}}")
	})

	assert.NoError(t, err)
	assert.Equal(t, "Foo is: foo Bar is: bar\n", stdout)
	assert.Empty(t, stderr)
}

func TestDescribeJSON(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return util.DescribeJSON(struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}{
			Foo: "foo",
			Bar: "bar",
		})
	})

	assert.NoError(t, err)
	assert.JSONEq(t, `{"foo":"foo", "bar": "bar"}`, stdout)
	assert.Empty(t, stderr)
}

func TestDescribeYAML(t *testing.T) {
	stdout, stderr, err := testutil.CaptureOutStreams(func() error {
		return util.DescribeYAML(struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		}{
			Foo: "foo",
			Bar: "bar",
		})
	})

	assert.NoError(t, err)
	assert.YAMLEq(t, `{"foo":"foo", "bar": "bar"}`, stdout)
	assert.Empty(t, stderr)
}

func TestWrap(t *testing.T) {
	wrapped := util.Wrap("json", map[string]interface{}{
		"foo": "bar",
	})
	jsonString, _ := json.Marshal(wrapped)
	assert.JSONEq(t, `{"json": {"foo": "bar"}}`, string(jsonString))
}

func TestValidateRequiredFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("foo", "", "foo")
	flags.String("bar", "", "bar")
	flags.String("baz", "", "baz")
	_ = flags.Set("foo", "foo")
	_ = flags.Set("bar", "bar")

	err := util.ValidateRequiredFlags(flags, "foo")
	assert.NoError(t, err)

	err = util.ValidateRequiredFlags(flags, "baz")
	assert.EqualError(t, err, "hcloud: required flag(s) \"baz\" not set")
}

func TestAddGroup(t *testing.T) {
	cmd := &cobra.Command{}
	util.AddGroup(cmd, "id", "title", &cobra.Command{})
	assert.Equal(t, len(cmd.Commands()), 1)
	assert.Equal(t, len(cmd.Groups()), 1)
	assert.Equal(t, cmd.Groups()[0].ID, "id")
	assert.Equal(t, cmd.Groups()[0].Title, "title:")
}

func TestToKebabCase(t *testing.T) {
	assert.Equal(t, "foo-bar", util.ToKebabCase("Foo Bar"))
	assert.Equal(t, "foo", util.ToKebabCase("Foo"))
}
