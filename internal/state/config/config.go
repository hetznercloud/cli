package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Config interface {
	// Write writes the config to the given writer. If w is nil, the config is written to the config file.
	Write(w io.Writer) error

	ParseConfig() error

	ActiveContext() Context
	SetActiveContext(Context)
	Contexts() []Context
	SetContexts([]Context)

	Preferences() Preferences
}

type schema struct {
	ActiveContext string      `toml:"active_context"`
	Preferences   preferences `toml:"preferences"`
	Contexts      []*context  `toml:"contexts"`
}

type config struct {
	path          string
	activeContext *context
	contexts      []*context
	preferences   preferences
}

var FlagSet *pflag.FlagSet

func init() {
	ResetFlags()
}

func ResetFlags() {
	FlagSet = pflag.NewFlagSet("hcloud", pflag.ContinueOnError)
	for _, o := range opts {
		o.AddToFlagSet(FlagSet)
	}
	if err := viper.BindPFlags(FlagSet); err != nil {
		panic(err)
	}
}

func NewConfig() Config {
	return &config{}
}

func ReadConfig(cfg Config) error {

	viper.SetConfigType("toml")
	viper.SetEnvPrefix("HCLOUD")

	// error is ignored since invalid flags are already handled by cobra
	_ = FlagSet.Parse(os.Args[1:])

	// load env already so we can determine the active context
	viper.AutomaticEnv()

	// load active context
	if err := cfg.ParseConfig(); err != nil {
		return err
	}

	return nil
}

func (cfg *config) ParseConfig() error {
	var s schema

	cfg.path = OptionConfig.Value()

	// read config file
	cfgBytes, err := os.ReadFile(cfg.path)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(cfgBytes, &s); err != nil {
		return err
	}

	// read config file into viper (particularly active_context)
	if err := viper.ReadConfig(bytes.NewReader(cfgBytes)); err != nil {
		return err
	}

	// read active context from viper
	if ctx := OptionContext.Value(); ctx != "" {
		s.ActiveContext = ctx
	}

	cfg.contexts = s.Contexts
	for i, ctx := range s.Contexts {
		if ctx.ContextName == s.ActiveContext {
			cfg.activeContext = cfg.contexts[i]
		}
	}

	if s.ActiveContext != "" && cfg.activeContext == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: active context %q not found\n", s.ActiveContext)
	}

	// load global preferences first so that contexts can override them
	if err = cfg.loadPreferences(cfg.preferences); err != nil {
		return err
	}

	// load context preferences
	if cfg.activeContext != nil {
		if err = cfg.loadPreferences(cfg.activeContext.ContextPreferences); err != nil {
			return err
		}
		// read context into viper (particularly the token)
		ctxBytes, err := toml.Marshal(cfg.activeContext)
		if err != nil {
			return err
		}
		if err = viper.ReadConfig(bytes.NewReader(ctxBytes)); err != nil {
			return err
		}
	}
	return nil
}

func (cfg *config) loadPreferences(prefs preferences) error {
	if err := prefs.validate(); err != nil {
		return err
	}
	ctxBytes, err := toml.Marshal(prefs)
	if err != nil {
		return err
	}
	return viper.MergeConfig(bytes.NewReader(ctxBytes))
}

func addOption[T any](flagFunc func(string, T, string) *T, key string, defaultVal T, usage string) {
	if flagFunc != nil {
		flagFunc(key, defaultVal, usage)
	}
	viper.SetDefault(key, defaultVal)
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

	var activeContext string
	if cfg.activeContext != nil {
		activeContext = cfg.activeContext.ContextName
	}

	s := schema{
		ActiveContext: activeContext,
		Preferences:   cfg.preferences,
		Contexts:      cfg.contexts,
	}

	return toml.NewEncoder(w).Encode(s)
}

func (cfg *config) ActiveContext() Context {
	return cfg.activeContext
}

func (cfg *config) SetActiveContext(ctx Context) {
	if ctx, ok := ctx.(*context); !ok {
		panic("invalid context type")
	} else {
		cfg.activeContext = ctx
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
}

func (cfg *config) Preferences() Preferences {
	if cfg.preferences == nil {
		cfg.preferences = make(preferences)
	}
	return cfg.preferences
}

func GetHcloudOpts(cfg Config) []hcloud.ClientOption {
	var opts []hcloud.ClientOption

	token := OptionToken.Value()

	opts = append(opts, hcloud.WithToken(token))
	if ep := OptionEndpoint.Value(); ep != "" {
		opts = append(opts, hcloud.WithEndpoint(ep))
	}
	if OptionDebug.Value() {
		if filePath := OptionDebugFile.Value(); filePath == "" {
			opts = append(opts, hcloud.WithDebugWriter(os.Stderr))
		} else {
			writer, _ := os.Create(filePath)
			opts = append(opts, hcloud.WithDebugWriter(writer))
		}
	}
	pollInterval := OptionPollInterval.Value()
	if pollInterval > 0 {
		opts = append(opts, hcloud.WithBackoffFunc(hcloud.ConstantBackoff(pollInterval)))
	}

	return opts
}
