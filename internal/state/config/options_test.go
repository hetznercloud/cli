package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	for _, opt := range Options {
		kind := reflect.TypeOf(opt.T()).Kind()
		if kind == reflect.Slice && !opt.HasFlags(OptionFlagSlice) {
			t.Errorf("option %s is a slice but does not have the slice flag", opt.GetName())
		}
		if kind != reflect.Slice && opt.HasFlags(OptionFlagSlice) {
			t.Errorf("option %s is not a slice but has the slice flag", opt.GetName())
		}
		if opt.HasFlags(OptionFlagPFlag | OptionFlagSensitive) {
			t.Errorf("%s: sensitive options shouldn't have pflags", opt.GetName())
		}
	}
}

func TestOption_HasFlags(t *testing.T) {
	opt := &Option[any]{Flags: OptionFlagSensitive | OptionFlagPFlag | OptionFlagSlice}
	assert.True(t, opt.HasFlags(OptionFlagSensitive))
	assert.True(t, opt.HasFlags(OptionFlagPFlag))
	assert.True(t, opt.HasFlags(OptionFlagSlice))
	assert.True(t, opt.HasFlags(OptionFlagSensitive|OptionFlagPFlag))
	assert.True(t, opt.HasFlags(OptionFlagSensitive|OptionFlagSlice))
	assert.True(t, opt.HasFlags(OptionFlagPFlag|OptionFlagSlice))
	assert.True(t, opt.HasFlags(OptionFlagSensitive|OptionFlagPFlag|OptionFlagSlice))
	assert.False(t, opt.HasFlags(OptionFlagConfig))
	assert.False(t, opt.HasFlags(OptionFlagConfig|OptionFlagSensitive))
}

func TestOption_EnvVar(t *testing.T) {
	opt := &Option[any]{Name: "foo", Flags: OptionFlagEnv}
	assert.Equal(t, "HCLOUD_FOO", opt.EnvVar())
	opt.Name = "foo-bar"
	assert.Equal(t, "HCLOUD_FOO_BAR", opt.EnvVar())
	opt.Name = "foo.bar-baz"
	assert.Equal(t, "HCLOUD_FOO_BAR_BAZ", opt.EnvVar())
	opt.Flags = 0
	assert.Empty(t, opt.EnvVar())
}
