package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/hetznercloud/cli/internal/cmd/util"
)

type Config interface {
	// Write writes the config to the given writer. If w is nil, the config is written to the config file.
	Write(w io.Writer) error

	// LoadDefault reads flags, environment variables, and the configured default file.
	LoadDefault() error
	// LoadFile reads flags, environment variables, and the specified file.
	LoadFile(path string) error
	// LoadReader reads flags, environment variables, and TOML from r.
	LoadReader(r io.Reader) error

	// ActiveContext returns the currently active context
	ActiveContext() Context
	// SetActiveContext sets the currently active context and also modifies the schema to reflect this change
	// This does NOT change any configuration values. Use [config.Config.UseConfig] to read the actual context into memory.
	SetActiveContext(Context) error
	// Contexts returns a list of currently loaded contexts
	Contexts() []Context
	// SetContexts sets the list of contexts and also modifies the schema to reflect this change
	SetContexts([]Context) error
	// UseContext temporarily switches context to the given context name and reloads the config, loading the values of the given context.
	// If name is nil, the context is unloaded and only the global preferences are used.
	// This change will not be written to the schema, so `active_context` will not be changed after writing.
	UseContext(name *string) error

	// Preferences returns the global preferences (as opposed to [Context.Preferences])
	Preferences() Preferences
	// Viper returns the currently active instance of viper
	Viper() *viper.Viper
	// FlagSet returns the FlagSet that options are bound to
	FlagSet() *pflag.FlagSet
	Options() []IOption
	LookupOption(name string) (IOption, bool)

	// Path returns the path to the used config file
	Path() string
	// Schema returns the TOML schema of the config file as a struct
	Schema() *Schema
}

type Schema struct {
	ActiveContext string      `toml:"active_context"`
	Preferences   Preferences `toml:"preferences"`
	Contexts      []*context  `toml:"contexts,omitempty"`
}

type config struct {
	v             *viper.Viper
	fs            *pflag.FlagSet
	args          []string
	warnings      io.Writer
	defaultPath   string
	initErr       error
	configBytes   []byte
	contextForced bool
	forcedContext string
	options       []IOption
	optionByName  map[string]IOption
	path          string
	activeContext *context
	contexts      []*context
	preferences   Preferences
	schema        Schema
}

type NewOption func(*config)

func WithArgs(args []string) NewOption {
	return func(cfg *config) {
		cfg.args = slices.Clone(args)
	}
}

func WithWarningWriter(w io.Writer) NewOption {
	return func(cfg *config) {
		cfg.warnings = w
	}
}

func WithDefaultPath(path string) NewOption {
	return func(cfg *config) {
		cfg.defaultPath = path
	}
}

func WithOptions(options ...IOption) NewOption {
	return func(cfg *config) {
		cfg.options = append(cfg.options, options...)
	}
}

func New(opts ...NewOption) Config {
	cfg := &config{warnings: io.Discard, options: DefaultOptions()}
	for _, opt := range opts {
		opt(cfg)
	}
	cfg.normalizeOptions()
	cfg.reset()
	return cfg
}

func (cfg *config) normalizeOptions() {
	options := make([]IOption, 0, len(cfg.options))
	indexByName := make(map[string]int, len(cfg.options))
	for _, option := range cfg.options {
		if index, ok := indexByName[option.GetName()]; ok {
			options[index] = option
			continue
		}
		indexByName[option.GetName()] = len(options)
		options = append(options, option)
	}
	cfg.options = options
	cfg.optionByName = make(map[string]IOption, len(options))
	for _, option := range options {
		cfg.optionByName[option.GetName()] = option
	}
}

func (cfg *config) reset() {
	cfg.initErr = nil
	cfg.v = viper.New()
	cfg.v.SetConfigType("toml")
	cfg.v.SetEnvPrefix("HCLOUD")
	cfg.v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	cfg.fs = pflag.NewFlagSet("hcloud", pflag.ContinueOnError)
	cfg.fs.Usage = func() {} // disable usage output
	for _, o := range cfg.options {
		cfg.initErr = errors.Join(cfg.initErr, o.addToFlagSet(cfg.fs))
		cfg.initErr = errors.Join(cfg.initErr, cfg.v.BindEnv(o.GetName(), o.EnvVar()))
	}
	cfg.initErr = errors.Join(cfg.initErr, cfg.v.BindPFlags(cfg.fs))
}

func (cfg *config) LoadDefault() error {
	return cfg.loadFile("")
}

func (cfg *config) LoadFile(path string) error {
	return cfg.loadFile(path)
}

func (cfg *config) LoadReader(r io.Reader) error {
	if err := cfg.prepare(); err != nil {
		return err
	}
	cfgBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return cfg.apply(cfgBytes)
}

func (cfg *config) prepare() error {
	cfg.schema = Schema{}
	cfg.activeContext = nil
	cfg.contexts = nil
	cfg.preferences = nil
	cfg.reset()
	if cfg.initErr != nil {
		return cfg.initErr
	}

	// we ignore unknown flags since we are only interested in config option flags
	cfg.fs.ParseErrorsAllowlist.UnknownFlags = true

	err := cfg.fs.Parse(cfg.args)
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		return err
	}

	cfg.path, err = OptionConfig.Get(cfg)
	if err != nil {
		return err
	}

	if cfg.path == "" {
		if cfg.defaultPath != "" {
			cfg.path = cfg.defaultPath
		} else {
			cfg.path = DefaultConfigPath()
		}
	}
	return nil
}

func (cfg *config) loadFile(path string) error {
	if err := cfg.prepare(); err != nil {
		return err
	}
	if path != "" {
		cfg.path = path
	}
	cfgBytes, err := os.ReadFile(cfg.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return cfg.apply(cfgBytes)
}

func (cfg *config) apply(cfgBytes []byte) error {
	cfg.configBytes = bytes.Clone(cfgBytes)
	if err := toml.Unmarshal(cfg.configBytes, &cfg.schema); err != nil {
		return err
	}

	if !cfg.contextForced && cfg.schema.ActiveContext != "" {
		// ReadConfig resets the current config and reads the new values
		// We don't use viper.Set here because of the value hierarchy. We want the env and flags to
		// be able to override the currently active context. viper.Set would take precedence over
		// env and flags.
		if err := cfg.v.ReadConfig(bytes.NewReader(fmt.Appendf(nil, "context = %q\n", cfg.schema.ActiveContext))); err != nil {
			return err
		}
	}

	activeContext := cfg.forcedContext
	if !cfg.contextForced {
		var err error
		activeContext, err = OptionContext.Get(cfg)
		if err != nil {
			return err
		}
	}

	cfg.contexts = cfg.schema.Contexts
	for i, ctx := range cfg.contexts {
		if ctx.ContextName == activeContext {
			cfg.activeContext = cfg.contexts[i]
			break
		}
	}

	if activeContext != "" && cfg.activeContext == nil {
		if _, err := fmt.Fprintf(cfg.warnings, "Warning: active context %q not found\n", activeContext); err != nil {
			return err
		}
	}

	// merge global preferences first so that contexts can override them
	cfg.preferences = cfg.schema.Preferences
	if err := cfg.preferences.merge(cfg.v, cfg.optionByName); err != nil {
		return err
	}

	if cfg.activeContext != nil {
		// Merge preferences into viper
		if err := cfg.activeContext.ContextPreferences.merge(cfg.v, cfg.optionByName); err != nil {
			return err
		}
		// Merge token into viper
		// We use viper.MergeConfig here for the same reason as above, except for
		// that we merge the config instead of replacing it.
		if err := cfg.v.MergeConfig(bytes.NewReader(fmt.Appendf(nil, `token = "%s"`, cfg.activeContext.ContextToken))); err != nil {
			return err
		}
	}
	return nil
}

func (cfg *config) Write(w io.Writer) error {
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

	var contents bytes.Buffer
	if err := toml.NewEncoder(&contents).Encode(s); err != nil {
		return err
	}

	if w != nil {
		if err := writeAll(w, contents.Bytes()); err != nil {
			return err
		}
		cfg.configBytes = bytes.Clone(contents.Bytes())
		return nil
	}

	if err := writeFileAtomic(cfg.path, contents.Bytes()); err != nil {
		return err
	}
	cfg.configBytes = bytes.Clone(contents.Bytes())
	return nil
}

func writeAll(w io.Writer, contents []byte) error {
	written, err := w.Write(contents)
	if err != nil {
		return err
	}
	if written != len(contents) {
		return io.ErrShortWrite
	}
	return nil
}

func writeFileAtomic(path string, contents []byte) (err error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	temporary, err := os.CreateTemp(dir, ".hcloud-config-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	committed := false
	defer func() {
		if closeErr := temporary.Close(); closeErr != nil && !errors.Is(closeErr, os.ErrClosed) {
			err = errors.Join(err, closeErr)
		}
		if !committed {
			if removeErr := os.Remove(temporaryPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				err = errors.Join(err, removeErr)
			}
		}
	}()

	if err := temporary.Chmod(0600); err != nil {
		return err
	}
	if err := writeAll(temporary, contents); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	committed = true
	return nil
}

func (cfg *config) ActiveContext() Context {
	return cfg.activeContext
}

func (cfg *config) SetActiveContext(ctx Context) error {
	if util.IsNil(ctx) {
		cfg.activeContext = nil
		cfg.schema.ActiveContext = ""
		return nil
	}
	ctxStruct, ok := ctx.(*context)
	if !ok {
		return fmt.Errorf("unsupported context implementation %T", ctx)
	}
	cfg.activeContext = ctxStruct
	cfg.schema.ActiveContext = ctxStruct.ContextName
	return nil
}

func (cfg *config) Contexts() []Context {
	ctxs := make([]Context, 0, len(cfg.contexts))
	for _, c := range cfg.contexts {
		ctxs = append(ctxs, c)
	}
	return ctxs
}

func (cfg *config) SetContexts(contexts []Context) error {
	converted := make([]*context, 0, len(contexts))
	for _, c := range contexts {
		c, ok := c.(*context)
		if !ok {
			return fmt.Errorf("unsupported context implementation %T", c)
		}
		converted = append(converted, c)
	}
	cfg.contexts = converted
	cfg.schema.Contexts = cfg.contexts
	return nil
}

func (cfg *config) UseContext(name *string) error {
	cfg.contextForced = true
	cfg.forcedContext = ""
	if name != nil {
		cfg.forcedContext = *name
	}

	contents := bytes.Clone(cfg.configBytes)
	if err := cfg.prepare(); err != nil {
		return err
	}
	return cfg.apply(contents)
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

func (cfg *config) Options() []IOption {
	return slices.Clone(cfg.options)
}

func (cfg *config) LookupOption(name string) (IOption, bool) {
	option, ok := cfg.optionByName[name]
	return option, ok
}

func (cfg *config) Path() string {
	return cfg.path
}

func (cfg *config) Schema() *Schema {
	return &cfg.schema
}
