package config

import (
	"fmt"
	"strings"
	"time"
)

type Preferences interface {
	Set(key string, value string) error
}

// preferences are options that can be set in the config file, globally or per context
type preferences map[string]any

func (p preferences) validate() error {
	for key := range p {
		opt, ok := opts[key]
		if !ok || !opt.HasSource(OptionSourcePreference) {
			return fmt.Errorf("unknown preference: %s", key)
		}
	}
	return nil
}

func (p preferences) Set(key string, value string) error {
	opt, ok := opts[key]
	if !ok || !opt.HasSource(OptionSourcePreference) {
		return fmt.Errorf("unknown preference: %s", key)
	}

	var val any
	switch t := opt.T().(type) {
	case bool:
		switch strings.ToLower(value) {
		case "true", "t", "yes", "y", "1":
			val = true
		case "false", "f", "no", "n", "0":
			val = false
		default:
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case string:
		val = value
	case time.Duration:
		var err error
		val, err = time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration value: %s", value)
		}
	case []string:
		val = strings.Split(value, ",")
	default:
		return fmt.Errorf("unsupported type %T", t)
	}

	configKey := strings.ReplaceAll(strings.ToLower(key), "-", "_")

	p[configKey] = val
	return nil
}

var _ Preferences = preferences{}
