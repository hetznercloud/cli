package config

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"

	"github.com/hetznercloud/cli/internal/cmd/util"
)

// Preferences are options that can be set in the config file, globally or per context
type Preferences map[string]any

func (p Preferences) Set(key string, values []string) (any, error) {
	opt, ok := Options[key]
	if !ok || !opt.HasFlag(OptionFlagPreference) {
		return nil, fmt.Errorf("unknown preference: %s", key)
	}

	var val any
	switch t := opt.T().(type) {
	case bool:
		if len(values) != 1 {
			return nil, fmt.Errorf("expected exactly one value")
		}
		value := values[0]
		switch strings.ToLower(value) {
		case "true", "t", "yes", "y", "1":
			val = true
		case "false", "f", "no", "n", "0":
			val = false
		default:
			return nil, fmt.Errorf("invalid boolean value: %s", value)
		}
	case string:
		if len(values) != 1 {
			return nil, fmt.Errorf("expected exactly one value")
		}
		val = values[0]
	case time.Duration:
		if len(values) != 1 {
			return nil, fmt.Errorf("expected exactly one value")
		}
		value := values[0]
		var err error
		val, err = time.ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("invalid duration value: %s", value)
		}
	case []string:
		newVal := values[:]
		slices.Sort(newVal)
		newVal = slices.Compact(newVal)
		val = newVal
	default:
		return nil, fmt.Errorf("unsupported type %T", t)
	}

	configKey := strings.ReplaceAll(strings.ToLower(key), "-", "_")

	p[configKey] = val
	return val, nil
}

func (p Preferences) Unset(key string) (bool, error) {
	opt, ok := Options[key]
	if !ok || !opt.HasFlag(OptionFlagPreference) {
		return false, fmt.Errorf("unknown preference: %s", key)
	}

	configKey := strings.ReplaceAll(strings.ToLower(key), "-", "_")
	_, ok = p[configKey]
	delete(p, configKey)
	return ok, nil
}

func (p Preferences) Add(key string, values []string) ([]any, error) {
	opt, ok := Options[key]
	if !ok || !opt.HasFlag(OptionFlagPreference) {
		return nil, fmt.Errorf("unknown preference: %s", key)
	}

	var added []any

	configKey := strings.ReplaceAll(strings.ToLower(key), "-", "_")
	val := p[configKey]
	switch opt.T().(type) {
	case []string:
		before := util.AnyToStringSlice(val)
		newVal := append(before, values...)
		slices.Sort(newVal)
		newVal = slices.Compact(newVal)
		val = newVal
		added = util.ToAnySlice(util.SliceDiff[[]string](newVal, before))
	default:
		return nil, fmt.Errorf("%s is not a list", key)
	}

	p[configKey] = val
	return added, nil
}

func (p Preferences) Remove(key string, values []string) ([]any, error) {
	opt, ok := Options[key]
	if !ok || !opt.HasFlag(OptionFlagPreference) {
		return nil, fmt.Errorf("unknown preference: %s", key)
	}

	var removed []any

	configKey := strings.ReplaceAll(strings.ToLower(key), "-", "_")
	val := p[configKey]
	switch opt.T().(type) {
	case []string:
		before := util.AnyToStringSlice(val)
		diff := util.SliceDiff[[]string](before, values)
		val = diff
		removed = util.ToAnySlice(util.SliceDiff[[]string](before, diff))
	default:
		return nil, fmt.Errorf("%s is not a list", key)
	}

	if reflect.ValueOf(val).Len() == 0 {
		delete(p, configKey)
	} else {
		p[configKey] = val
	}
	return removed, nil
}

func (p Preferences) merge(v *viper.Viper) error {
	if err := p.validate(); err != nil {
		return err
	}
	m := make(map[string]any)
	for k, v := range p {
		m[strings.ReplaceAll(k, "_", "-")] = v
	}
	var buf bytes.Buffer
	err := toml.NewEncoder(&buf).Encode(m)
	if err != nil {
		return err
	}
	return v.MergeConfig(&buf)
}

func (p Preferences) validate() error {
	for key := range p {
		opt, ok := Options[strings.ReplaceAll(key, "_", "-")]
		if !ok || !opt.HasFlag(OptionFlagPreference) {
			return fmt.Errorf("unknown preference: %s", key)
		}
	}
	return nil
}

var _ Preferences = Preferences{}
