package state

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

var DefaultConfigPath string

type Config struct {
	Endpoint      string
	ActiveContext *ConfigContext
	Contexts      []*ConfigContext
}

type ConfigContext struct {
	Name  string
	Token string
}

func (config *Config) ContextNames() []string {
	if len(config.Contexts) == 0 {
		return nil
	}
	names := make([]string, len(config.Contexts))
	for i, ctx := range config.Contexts {
		names[i] = ctx.Name
	}
	return names
}

func (config *Config) ContextByName(name string) *ConfigContext {
	for _, c := range config.Contexts {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (config *Config) RemoveContext(context *ConfigContext) {
	for i, c := range config.Contexts {
		if c == context {
			config.Contexts = append(config.Contexts[:i], config.Contexts[i+1:]...)
			return
		}
	}
}

type RawConfig struct {
	ActiveContext string             `toml:"active_context,omitempty"`
	Contexts      []RawConfigContext `toml:"contexts"`
}

type RawConfigContext struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}

func MarshalConfig(c *Config) ([]byte, error) {
	if c == nil {
		return []byte{}, nil
	}

	var raw RawConfig
	if c.ActiveContext != nil {
		raw.ActiveContext = c.ActiveContext.Name
	}
	for _, context := range c.Contexts {
		raw.Contexts = append(raw.Contexts, RawConfigContext{
			Name:  context.Name,
			Token: context.Token,
		})
	}
	return toml.Marshal(raw)
}

func UnmarshalConfig(config *Config, data []byte) error {
	var raw RawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, rawContext := range raw.Contexts {
		config.Contexts = append(config.Contexts, &ConfigContext{
			Name:  rawContext.Name,
			Token: rawContext.Token,
		})
	}
	if raw.ActiveContext != "" {
		for _, c := range config.Contexts {
			if c.Name == raw.ActiveContext {
				config.ActiveContext = c
				break
			}
		}
		if config.ActiveContext == nil {
			return fmt.Errorf("active context %s not found", raw.ActiveContext)
		}
	}
	return nil
}
