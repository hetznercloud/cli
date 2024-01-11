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
}

type Context struct {
	Name    string
	Token   string
	SSHKeys []string
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

func ContextNames(cfg Config) []string {
	ctxs := cfg.Contexts()
	names := make([]string, len(ctxs))
	for i, ctx := range ctxs {
		names[i] = ctx.Name
	}
	return names
}

func ContextByName(cfg Config, name string) *Context {
	for _, c := range cfg.Contexts() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func RemoveContext(cfg Config, context *Context) {
	var filtered []*Context
	for _, c := range cfg.Contexts() {
		if c != context {
			filtered = append(filtered, c)
		}
	}
	cfg.SetContexts(filtered)
}

type rawConfig struct {
	ActiveContext string             `toml:"active_context,omitempty"`
	Contexts      []rawConfigContext `toml:"contexts"`
}

type rawConfigContext struct {
	Name    string   `toml:"name"`
	Token   string   `toml:"token"`
	SSHKeys []string `toml:"ssh_keys"`
}

func (cfg *config) marshal() ([]byte, error) {
	var raw rawConfig
	if cfg.activeContext != nil {
		raw.ActiveContext = cfg.activeContext.Name
	}
	for _, context := range cfg.contexts {
		raw.Contexts = append(raw.Contexts, rawConfigContext{
			Name:    context.Name,
			Token:   context.Token,
			SSHKeys: context.SSHKeys,
		})
	}
	return toml.Marshal(raw)
}

func (cfg *config) unmarshal(data []byte) error {
	var raw rawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, rawContext := range raw.Contexts {
		cfg.contexts = append(cfg.contexts, &Context{
			Name:    rawContext.Name,
			Token:   rawContext.Token,
			SSHKeys: rawContext.SSHKeys,
		})
	}
	if raw.ActiveContext != "" {
		for _, c := range cfg.contexts {
			if c.Name == raw.ActiveContext {
				cfg.activeContext = c
				break
			}
		}
		if cfg.activeContext == nil {
			return fmt.Errorf("active context %s not found", raw.ActiveContext)
		}
	}
	return nil
}
