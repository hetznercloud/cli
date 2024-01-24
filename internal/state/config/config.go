package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//go:generate mockgen -package config -destination zz_config_mock.go github.com/hetznercloud/cli/internal/state/config Config

type Config interface {
	Write() error

	// LoadActiveContext loads values from the active context
	LoadActiveContext() error

	ActiveContext() *Context
	SetActiveContext(*Context)
	Contexts() []*Context
	SetContexts([]*Context)

	FlagSet() *pflag.FlagSet

	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
}

type Context struct {
	Name  string
	Token string
}

type config struct {
	flagSet *pflag.FlagSet
	Config
}

func ReadConfig() (Config, error) {

	cfg := &config{
		flagSet: pflag.NewFlagSet("hcloud", pflag.ContinueOnError),
	}

	viper.SetConfigType("toml")
	viper.SetEnvPrefix("HCLOUD")

	addFlag(cfg.flagSet.String, "config", DefaultConfigPath(), "Config file path")
	addFlag(cfg.flagSet.String, "token", "", "Hetzner Cloud API token")
	addFlag(cfg.flagSet.String, "endpoint", "", "Hetzner Cloud API endpoint")
	addFlag(cfg.flagSet.Bool, "debug", false, "Enable debug output")
	addFlag(cfg.flagSet.String, "debug-file", "", "Write debug output to file")
	addFlag(cfg.flagSet.String, "context", "", "Active context")
	addFlag(cfg.flagSet.Duration, "poll-interval", 500*time.Millisecond, "Interval at which to poll information, for example action progress")
	addFlag(cfg.flagSet.Bool, "quiet", false, "Only print error messages")
	addFlag(cfg.flagSet.StringSlice, "default-ssh-keys", nil, "Default SSH keys for new servers")

	if err := cfg.flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	if err := viper.BindPFlags(cfg.flagSet); err != nil {
		return nil, err
	}

	// load env already so we can determine the active context
	viper.AutomaticEnv()

	// load active context
	if err := cfg.LoadActiveContext(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (*config) LoadActiveContext() error {

	var schema struct {
		ActiveContext string           `toml:"active_context"`
		Contexts      []map[string]any `toml:"contexts"`
	}

	// read config file
	cfgBytes, err := os.ReadFile(viper.GetString("config"))
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(cfgBytes, &schema); err != nil {
		return err
	}

	// read config file into viper (particularly active_context)
	if err := viper.ReadConfig(bytes.NewReader(cfgBytes)); err != nil {
		return err
	}

	// read active context from viper
	if ctx := viper.GetString("context"); ctx != "" {
		schema.ActiveContext = ctx
	}

	// find active context in schema
	var activeContext any
	for _, ctx := range schema.Contexts {
		if name, ok := ctx["name"].(string); ok && name == schema.ActiveContext {
			activeContext = ctx
			break
		}
	}

	if activeContext == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: active context %q not found\n", schema.ActiveContext)
	}

	// merge active context into viper
	ctxBytes, err := toml.Marshal(activeContext)
	if err := viper.MergeConfig(bytes.NewReader(ctxBytes)); err != nil {
		return err
	}
	return nil
}

func addFlag[T any](f func(string, T, string) *T, key string, defaultVal T, usage string) {
	f(key, defaultVal, usage)
	viper.SetDefault(key, defaultVal)
}

func (cfg *config) FlagSet() *pflag.FlagSet {
	return cfg.flagSet
}

func (*config) Get(key string) any {
	return viper.Get(key)
}

func (*config) GetString(key string) string {
	return viper.GetString(key)
}

func (*config) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (*config) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func (*config) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (cfg *config) Write() error {
	return fmt.Errorf("not implemented")
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
