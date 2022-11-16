package state

import (
	"fmt"

	toml "github.com/pelletier/go-toml/v2"
)

var DefaultConfigPath string

type Config struct {
	Endpoint           string
	ActiveContext      *ConfigContext
	Contexts           []*ConfigContext
	SubcommandDefaults map[string]*SubcommandDefaults
}

type ConfigContext struct {
	Name  string
	Token string
}

type SubcommandDefaults struct {
	Sorting        []string
	DefaultColumns []string
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
	ActiveContext      string                            `toml:"active_context,omitempty"`
	Contexts           []RawConfigContext                `toml:"contexts"`
	SubcommandDefaults map[string]*RAWSubcommandDefaults `toml:"defaults,omitempty"`
}

type RawConfigContext struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}

type RAWSubcommandDefaults struct {
	Sorting []string `toml:"sort,omitempty"`
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
	if len(c.SubcommandDefaults) != 0 {
		raw.SubcommandDefaults = make(map[string]*RAWSubcommandDefaults)

		for command, defaults := range c.SubcommandDefaults {
			raw.SubcommandDefaults[command] = &RAWSubcommandDefaults{
				Sorting: defaults.Sorting,
			}
		}
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
	if len(raw.SubcommandDefaults) > 0 {
		config.SubcommandDefaults = make(map[string]*SubcommandDefaults)
		for command, defaults := range raw.SubcommandDefaults {
			config.SubcommandDefaults[command] = &SubcommandDefaults{
				Sorting: defaults.Sorting,
			}
		}
	}

	return nil
}
