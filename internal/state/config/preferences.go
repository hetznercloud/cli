package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

// Preferences are options that can be set in the config file, globally or per context
type Preferences map[string]any

func (p Preferences) Get(key string) (any, bool) {
	var m map[string]any = p
	path := splitPath(key)
	for i, key := range path {
		if i == len(path)-1 {
			val, ok := m[key]
			return val, ok
		}
		next, ok := m[key].(map[string]any)
		if !ok {
			break
		}
		m = next
	}
	return nil, false
}

func (p Preferences) Set(key string, val any) {
	var m map[string]any = p
	path := splitPath(key)
	for i, key := range path {
		if i == len(path)-1 {
			m[key] = val
			return
		}
		if next, ok := m[key].(map[string]any); !ok {
			next = make(map[string]any)
			m[key] = next
			m = next
		} else {
			m = next
		}
	}
}

func (p Preferences) Unset(key string) bool {
	var m map[string]any = p
	path := splitPath(key)
	parents := make([]map[string]any, 0, len(path)-1)
	for i, key := range path {
		parents = append(parents, m)
		if i == len(path)-1 {
			_, ok := m[key]
			delete(m, key)
			// delete parent maps if they are empty
			for i := len(parents) - 1; i >= 0; i-- {
				if len(parents[i]) == 0 {
					if i > 0 {
						delete(parents[i-1], path[i-1])
					}
				}
			}
			return ok
		}
		next, ok := m[key].(map[string]any)
		if !ok {
			return false
		}
		m = next
	}
	return false
}

func (p Preferences) Validate() error {
	return validate(p, "")
}

func (p Preferences) merge(v *viper.Viper) error {
	if err := p.Validate(); err != nil {
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

func validate(m map[string]any, prefix string) error {
	for configKey, val := range m {
		key := prefix + strings.ReplaceAll(configKey, "_", "-")
		if val, ok := val.(map[string]any); ok {
			if err := validate(val, key+"."); err != nil {
				return err
			}
			continue
		}
		opt, ok := Options[key]
		if !ok || !opt.HasFlags(OptionFlagPreference) {
			return fmt.Errorf("unknown preference: %s", key)
		}
	}
	return nil
}

func splitPath(key string) []string {
	configKey := strings.ReplaceAll(strings.ToLower(key), "-", "_")
	return strings.Split(configKey, ".")
}
