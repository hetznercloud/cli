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

func TestIsNil(t *testing.T) {
	assert.True(t, util.IsNil(nil))
	assert.True(t, util.IsNil((*int)(nil)))
	assert.True(t, util.IsNil((chan int)(nil)))
	assert.True(t, util.IsNil((map[int]int)(nil)))
	assert.True(t, util.IsNil(([]int)(nil)))
	assert.True(t, util.IsNil((func())(nil)))
	assert.True(t, util.IsNil((interface{})(nil)))
	assert.True(t, util.IsNil((error)(nil)))
	assert.False(t, util.IsNil(0))
	assert.False(t, util.IsNil(""))
	assert.False(t, util.IsNil([]int{}))
	assert.False(t, util.IsNil(struct{}{}))
}

func TestFilterNil(t *testing.T) {
	type testStruct struct {
		a, b, c int //nolint:unused
	}
	assert.Equal(t, []interface{}{0, ""}, util.FilterNil([]interface{}{0, nil, ""}))
	assert.Equal(t, []*testStruct{{1, 2, 3}, {}}, util.FilterNil([]*testStruct{{1, 2, 3}, nil, {}, (*testStruct)(nil)}))
}

func TestSliceDiff(t *testing.T) {
	assert.Equal(t, []int{1, 2}, util.SliceDiff[[]int]([]int{1, 2, 3}, []int{3, 4}))
	assert.Equal(t, []int{4}, util.SliceDiff[[]int]([]int{3, 4}, []int{1, 2, 3}))
	assert.Empty(t, util.SliceDiff[[]int]([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.Empty(t, util.SliceDiff[[]int]([]int{}, []int{}))
	assert.Equal(t, []string{"a", "b"}, util.SliceDiff[[]string]([]string{"a", "b", "c"}, []string{"c", "d"}))
	assert.Equal(t, []string{"a"}, util.SliceDiff[[]string]([]string{"b", "a", "b", "b", "c", "c"}, []string{"b", "c"}))
}

func TestAnyToAnySlice(t *testing.T) {
	assert.Equal(t, []any{1, "foo", true}, util.AnyToAnySlice([]any{1, "foo", true}))
	assert.Equal(t, []any{"a", "b", "c"}, util.AnyToAnySlice([]string{"a", "b", "c"}))
	assert.Equal(t, []any{1, 2, 3}, util.AnyToAnySlice([]int{1, 2, 3}))
	assert.Equal(t, []any{true, false}, util.AnyToAnySlice([]bool{true, false}))
	assert.Nil(t, util.AnyToAnySlice(1))
	assert.Nil(t, util.AnyToAnySlice("abc"))
	assert.Nil(t, util.AnyToAnySlice(nil))
}

func TestAnyToStringSlice(t *testing.T) {
	assert.Equal(t, []string{"1", "foo", "true"}, util.AnyToStringSlice([]any{1, "foo", true}))
	assert.Equal(t, []string{"a", "b", "c"}, util.AnyToStringSlice([]string{"a", "b", "c"}))
	assert.Equal(t, []string{"1", "2", "3"}, util.AnyToStringSlice([]int{1, 2, 3}))
	assert.Equal(t, []string{"true", "false"}, util.AnyToStringSlice([]bool{true, false}))
	assert.Nil(t, util.AnyToStringSlice(1))
	assert.Nil(t, util.AnyToStringSlice("abc"))
	assert.Nil(t, util.AnyToStringSlice(nil))
}

func TestToStringSlice(t *testing.T) {
	assert.Equal(t, []string{"1", "foo", "true"}, util.ToStringSlice([]any{1, "foo", true}))
	assert.Equal(t, []string{"a", "b", "c"}, util.ToStringSlice([]any{"a", "b", "c"}))
	assert.Equal(t, []string{"1", "2", "3"}, util.ToStringSlice([]any{1, 2, 3}))
	assert.Equal(t, []string{"true", "false"}, util.ToStringSlice([]any{true, false}))
}

func TestToAnySlice(t *testing.T) {
	assert.Equal(t, []any{1, "foo", true}, util.ToAnySlice([]any{1, "foo", true}))
	assert.Equal(t, []any{"a", "b", "c"}, util.ToAnySlice([]string{"a", "b", "c"}))
	assert.Equal(t, []any{1, 2, 3}, util.ToAnySlice([]int{1, 2, 3}))
	assert.Equal(t, []any{true, false}, util.ToAnySlice([]bool{true, false}))
}

func TestParseBoolLenient(t *testing.T) {
	b, err := util.ParseBoolLenient("true")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ParseBoolLenient("True")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ParseBoolLenient("t")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ParseBoolLenient("yes")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ParseBoolLenient("y")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ParseBoolLenient("1")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ParseBoolLenient("false")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ParseBoolLenient("False")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ParseBoolLenient("f")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ParseBoolLenient("no")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ParseBoolLenient("n")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ParseBoolLenient("0")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ParseBoolLenient("invalid")
	assert.EqualError(t, err, "invalid boolean value: invalid")
	b, err = util.ParseBoolLenient("")
	assert.EqualError(t, err, "invalid boolean value: ")
}

func TestBoolFromAny(t *testing.T) {
	b, err := util.ToBoolE(true)
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ToBoolE("true")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ToBoolE("false")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ToBoolE("yes")
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ToBoolE("no")
	assert.NoError(t, err)
	assert.False(t, b)
	b, err = util.ToBoolE(1)
	assert.NoError(t, err)
	assert.True(t, b)
	b, err = util.ToBoolE(0)
	assert.NoError(t, err)
	assert.False(t, b)
	_, err = util.ToBoolE("invalid")
	assert.EqualError(t, err, "invalid boolean value: invalid")
}

func TestToStringSliceDelimited(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, util.ToStringSliceDelimited([]string{"a", "b", "c"}))
	assert.Equal(t, []string{"a", "b", "c"}, util.ToStringSliceDelimited([]any{"a", "b", "c"}))
	assert.Equal(t, []string{"a", "b", "c"}, util.ToStringSliceDelimited("a,b,c"))
	assert.Equal(t, []string{"0", "1", "2"}, util.ToStringSliceDelimited([]int{0, 1, 2}))
}

func TestRemoveDuplicates(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, util.RemoveDuplicates([]string{"a", "b", "c"}))
	assert.Equal(t, []string{"a", "b", "c"}, util.RemoveDuplicates([]string{"a", "b", "c", "a", "b", "c"}))
	assert.Equal(t, []string{"a", "b", "c"}, util.RemoveDuplicates([]string{"a", "b", "c", "c", "b", "a"}))
	assert.Equal(t, []string{"c", "b", "a"}, util.RemoveDuplicates([]string{"c", "b", "a", "a", "b", "c"}))
	assert.Equal(t, []string{"a"}, util.RemoveDuplicates([]string{"a", "a", "a", "a", "a"}))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, util.RemoveDuplicates([]int{1, 2, 1, 1, 3, 2, 1, 4, 3, 2, 5, 4, 3, 2, 1}))
}
