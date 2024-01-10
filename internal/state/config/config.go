package config

import (
	"fmt"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

//go:generate mockgen -package config -destination zz_config_mock.go github.com/hetznercloud/cli/internal/state/config Config

type Config interface {
	Write() error

	ActiveContext() *Context
	SetActiveContext(*Context)
	Contexts() []*Context
	SetContexts([]*Context)
	Endpoint() string
	SetEndpoint(string)

	ContextNames() []string
	ContextByName(name string) *Context
	RemoveContext(context *Context)
}

type Context struct {
	Name  string
	Token string
}

type config struct {
	path          string
	endpoint      string
	activeContext *Context   `toml:"active_context,omitempty"`
	contexts      []*Context `toml:"contexts"`
}

func ReadConfig(path string) (Config, error) {
	cfg := &config{path: path}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = cfg.unmarshal(data); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *config) Write() error {
	data, err := cfg.marshal()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cfg.path), 0777); err != nil {
		return err
	}
	if err := os.WriteFile(cfg.path, data, 0600); err != nil {
		return err
	}
	return nil
}

func (cfg *config) ActiveContext() *Context {
	return cfg.activeContext
}

func (cfg *config) SetActiveContext(context *Context) {
	cfg.activeContext = context
}

func (cfg *config) Contexts() []*Context {
	return cfg.contexts
}

func (cfg *config) SetContexts(contexts []*Context) {
	cfg.contexts = contexts
}

func (cfg *config) Endpoint() string {
	return cfg.endpoint
}

func (cfg *config) SetEndpoint(endpoint string) {
	cfg.endpoint = endpoint
}

func (cfg *config) ContextNames() []string {
	if len(cfg.contexts) == 0 {
		return nil
	}
	names := make([]string, len(cfg.contexts))
	for i, ctx := range cfg.contexts {
		names[i] = ctx.Name
	}
	return names
}

func (cfg *config) ContextByName(name string) *Context {
	for _, c := range cfg.contexts {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (cfg *config) RemoveContext(context *Context) {
	for i, c := range cfg.contexts {
		if c == context {
			cfg.contexts = append(cfg.contexts[:i], cfg.contexts[i+1:]...)
			return
		}
	}
}

type rawConfig struct {
	activeContext string
	contexts      []rawConfigContext
}

type rawConfigContext struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}

func (cfg *config) marshal() ([]byte, error) {
	var raw rawConfig
	if cfg.activeContext != nil {
		raw.activeContext = cfg.activeContext.Name
	}
	for _, context := range cfg.contexts {
		raw.contexts = append(raw.contexts, rawConfigContext{
			Name:  context.Name,
			Token: context.Token,
		})
	}
	return toml.Marshal(raw)
}

func (cfg *config) unmarshal(data []byte) error {
	var raw rawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, rawContext := range raw.contexts {
		cfg.contexts = append(cfg.contexts, &Context{
			Name:  rawContext.Name,
			Token: rawContext.Token,
		})
	}
	if raw.activeContext != "" {
		for _, c := range cfg.contexts {
			if c.Name == raw.activeContext {
				cfg.activeContext = c
				break
			}
		}
		if cfg.activeContext == nil {
			return fmt.Errorf("active context %s not found", raw.activeContext)
		}
	}
	return nil
}
