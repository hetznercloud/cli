package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/hetznercloud/cli/internal/cmd/util"
)

type Config interface {
	// Write writes the config to the given writer. If w is nil, the config is written to the config file.
	Write(w io.Writer) error

	Reset()
	ParseConfigFile(f any) error

	ActiveContext() Context
	SetActiveContext(Context)
	Contexts() []Context
	SetContexts([]Context)

	Preferences() Preferences
	Viper() *viper.Viper
	FlagSet() *pflag.FlagSet
	Path() string
	Schema() *Schema
}

type Schema struct {
	ActiveContext string      `toml:"active_context"`
	Preferences   Preferences `toml:"preferences"`
	Contexts      []*context  `toml:"contexts"`
}

type config struct {
	v             *viper.Viper
	fs            *pflag.FlagSet
	path          string
	activeContext *context
	contexts      []*context
	preferences   Preferences
	schema        Schema
}

func NewConfig() Config {
	cfg := &config{}
	cfg.Reset()
	return cfg
}

func (cfg *config) Reset() {
	cfg.v = viper.New()
	cfg.v.SetConfigType("toml")
	cfg.v.SetEnvPrefix("HCLOUD")
	cfg.v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cfg.fs = pflag.NewFlagSet("hcloud", pflag.ContinueOnError)
	for _, o := range Options {
		o.addToFlagSet(cfg.fs)
	}
	if err := cfg.v.BindPFlags(cfg.fs); err != nil {
		panic(err)
	}
}

// ReadConfig reads the config from the flags, env and the given config file f.
// See [ParseConfigFile] for the supported types of f.
func ReadConfig(cfg Config, f any) error {

	// error is ignored since invalid flags are already handled by cobra
	_ = cfg.FlagSet().Parse(os.Args[1:])

	// load env already so we can determine the active context
	cfg.Viper().AutomaticEnv()

	return cfg.ParseConfigFile(f)
}

// ParseConfigFile parses the given config file f.
// f can be of the following types:
// - nil: the default config file is used
// - string: the path to the config file
// - io.Reader: the config is read from the reader
// - []byte: the config is read from the byte slice
// - any other type: an error is returned
func (cfg *config) ParseConfigFile(f any) error {
	var (
		cfgBytes []byte
		err      error
	)

	cfg.path = OptionConfig.Get(cfg)
	path, ok := f.(string)
	if path != "" && ok {
		cfg.path = path
	}

	if f == nil || ok {
		// read config from file
		cfgBytes, err = os.ReadFile(cfg.path)
		if err != nil {
			return err
		}
	} else {
		switch f := f.(type) {
		case io.Reader:
			cfgBytes, err = io.ReadAll(f)
			if err != nil {
				return err
			}
		case []byte:
			cfgBytes = f
		default:
			return fmt.Errorf("invalid config file type %T", f)
		}
	}

	if err := toml.Unmarshal(cfgBytes, &cfg.schema); err != nil {
		return err
	}

	if cfg.schema.ActiveContext != "" {
		// ReadConfig resets the current config and reads the new values
		// We don't use viper.Set here because of the value hierarchy. We want the env and flags to
		// be able to override the currently active context. viper.Set would take precedence over
		// env and flags.
		err = cfg.v.ReadConfig(bytes.NewReader([]byte(fmt.Sprintf("context = %q\n", cfg.schema.ActiveContext))))
		if err != nil {
			return err
		}
	}

	// read active context from viper
	activeContext := cfg.schema.ActiveContext
	if ctx := OptionContext.Get(cfg); ctx != "" {
		activeContext = ctx
	}

	cfg.contexts = cfg.schema.Contexts
	for i, ctx := range cfg.contexts {
		if ctx.ContextName == activeContext {
			cfg.activeContext = cfg.contexts[i]
			break
		}
	}

	if cfg.schema.ActiveContext != "" && cfg.activeContext == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: active context %q not found\n", cfg.schema.ActiveContext)
	}

	// merge global preferences first so that contexts can override them
	cfg.preferences = cfg.schema.Preferences
	if err = cfg.preferences.merge(cfg.v); err != nil {
		return err
	}

	if cfg.activeContext != nil {
		// Merge preferences into viper
		if err = cfg.activeContext.ContextPreferences.merge(cfg.v); err != nil {
			return err
		}
		// Merge token into viper
		// We use viper.MergeConfig here for the same reason as above, except for
		// that we merge the config instead of replacing it.
		if err = cfg.v.MergeConfig(bytes.NewReader([]byte(fmt.Sprintf(`token = "%s"`, cfg.activeContext.ContextToken)))); err != nil {
			return err
		}
	}
	return nil
}

func (cfg *config) Write(w io.Writer) (err error) {
	if w == nil {
		f, err := os.OpenFile(cfg.path, os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer func() {
			err = errors.Join(err, f.Close())
		}()
		w = f
	}

	s := cfg.schema

	// this is so that we don't marshal empty preferences (this could happen e.g. after the last key is removed)
	if s.Preferences != nil && len(s.Preferences) == 0 {
		s.Preferences = nil
	}
	for _, ctx := range s.Contexts {
		if ctx.ContextPreferences != nil && len(ctx.ContextPreferences) == 0 {
			ctx.ContextPreferences = nil
		}
	}

	return toml.NewEncoder(w).Encode(s)
}

func (cfg *config) ActiveContext() Context {
	return cfg.activeContext
}

func (cfg *config) SetActiveContext(ctx Context) {
	if util.IsNil(ctx) {
		cfg.activeContext = nil
		cfg.schema.ActiveContext = ""
		return
	}
	if ctx, ok := ctx.(*context); !ok {
		panic("invalid context type")
	} else {
		cfg.activeContext = ctx
		cfg.schema.ActiveContext = ctx.ContextName
	}
}

func (cfg *config) Contexts() []Context {
	ctxs := make([]Context, 0, len(cfg.contexts))
	for _, c := range cfg.contexts {
		ctxs = append(ctxs, c)
	}
	return ctxs
}

func (cfg *config) SetContexts(contexts []Context) {
	cfg.contexts = make([]*context, 0, len(cfg.contexts))
	for _, c := range contexts {
		if c, ok := c.(*context); !ok {
			panic("invalid context type")
		} else {
			cfg.contexts = append(cfg.contexts, c)
		}
	}
	cfg.schema.Contexts = cfg.contexts
}

func (cfg *config) Preferences() Preferences {
	if cfg.preferences == nil {
		cfg.preferences = make(Preferences)
		cfg.schema.Preferences = cfg.preferences
	}
	return cfg.preferences
}

func (cfg *config) Viper() *viper.Viper {
	return cfg.v
}

func (cfg *config) FlagSet() *pflag.FlagSet {
	return cfg.fs
}

func (cfg *config) Path() string {
	return cfg.path
}

func (cfg *config) Schema() *Schema {
	return &cfg.schema
}
