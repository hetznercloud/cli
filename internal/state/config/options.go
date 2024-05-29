package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type OptionFlag int

const (
	// OptionFlagPreference indicates that the option can be set in the config file, globally or per context (in the preferences section)
	OptionFlagPreference OptionFlag = 1 << iota
	// OptionFlagConfig indicates that the option can be configured inside the configuration file
	OptionFlagConfig
	// OptionFlagPFlag indicates that the option can be set via a command line flag
	OptionFlagPFlag
	// OptionFlagEnv indicates that the option can be set via an environment variable
	OptionFlagEnv
	// OptionFlagSensitive indicates that the option holds sensitive data and should not be printed
	OptionFlagSensitive

	DefaultPreferenceFlags = OptionFlagPreference | OptionFlagConfig | OptionFlagPFlag | OptionFlagEnv
)

type IOption interface {
	addToFlagSet(fs *pflag.FlagSet)
	GetName() string
	GetDescription() string
	ConfigKey() string
	EnvVar() string
	FlagName() string
	HasFlag(src OptionFlag) bool
	GetAsAny(c Config) any
	OverrideAny(c Config, v any)
	Changed(c Config) bool
	Completions() []string
	IsSlice() bool
	T() any
}

type overrides struct {
	configKey string
	envVar    string
	flagName  string
}

var Options = make(map[string]IOption)

// Note: &^ is the bit clear operator and is used to remove flags from the default flag set
var (
	OptionConfig = newOpt(
		"config",
		"Config file path",
		DefaultConfigPath(),
		OptionFlagPFlag|OptionFlagEnv,
		nil,
	)

	OptionToken = newOpt(
		"token",
		"Hetzner Cloud API token",
		"",
		OptionFlagConfig|OptionFlagEnv|OptionFlagSensitive,
		nil,
	)

	OptionContext = newOpt(
		"context",
		"Currently active context",
		"",
		OptionFlagConfig|OptionFlagEnv|OptionFlagPFlag,
		&overrides{configKey: "active_context"},
	)

	OptionEndpoint = newOpt(
		"endpoint",
		"Hetzner Cloud API endpoint",
		hcloud.Endpoint,
		DefaultPreferenceFlags,
		nil,
	)

	OptionDebug = newOpt(
		"debug",
		"Enable debug output",
		false,
		DefaultPreferenceFlags,
		nil,
	)

	OptionDebugFile = newOpt(
		"debug-file",
		"File to write debug output to",
		"",
		DefaultPreferenceFlags,
		nil,
	)

	OptionPollInterval = newOpt(
		"poll-interval",
		"Interval at which to poll information, for example action progress",
		500*time.Millisecond,
		DefaultPreferenceFlags,
		nil,
	)

	OptionQuiet = newOpt(
		"quiet",
		"If true, only print error messages",
		false,
		DefaultPreferenceFlags,
		nil,
	)
)

type Option[T any] struct {
	Name        string
	Description string
	Default     T
	Flags       OptionFlag
	overrides   *overrides
}

func (o *Option[T]) Get(c Config) T {
	val := c.Viper().Get(o.Name)
	if val == nil {
		return o.Default
	}
	var t T
	switch any(t).(type) {
	case time.Duration:
		if v, ok := val.(string); ok {
			d, err := time.ParseDuration(v)
			if err != nil {
				panic(err)
			}
			val = d
		}
		if v, ok := val.(int64); ok {
			val = time.Duration(v)
		}
	case bool:
		if v, ok := val.(string); ok {
			b, err := strconv.ParseBool(v)
			if err != nil {
				panic(err)
			}
			val = b
		}
	case []string:
		if v, ok := val.([]any); ok {
			val = util.ToStringSlice(v)
		}
	}
	return val.(T)
}

func (o *Option[T]) GetAsAny(c Config) any {
	return o.Get(c)
}

func (o *Option[T]) Override(c Config, v T) {
	c.Viper().Set(o.Name, v)
}

func (o *Option[T]) OverrideAny(c Config, v any) {
	c.Viper().Set(o.Name, v)
}

func (o *Option[T]) Changed(c Config) bool {
	return c.Viper().IsSet(o.Name)
}

func (o *Option[T]) HasFlag(src OptionFlag) bool {
	return o.Flags&src != 0
}

func (o *Option[T]) IsSlice() bool {
	return reflect.TypeOf(o.T()).Kind() == reflect.Slice
}

func (o *Option[T]) GetName() string {
	return o.Name
}

func (o *Option[T]) GetDescription() string {
	return o.Description
}

func (o *Option[T]) ConfigKey() string {
	if !o.HasFlag(OptionFlagConfig) {
		return ""
	}
	if o.overrides != nil && o.overrides.configKey != "" {
		return o.overrides.configKey
	}
	return strings.ReplaceAll(strings.ToLower(o.Name), "-", "_")
}

func (o *Option[T]) EnvVar() string {
	if !o.HasFlag(OptionFlagEnv) {
		return ""
	}
	if o.overrides != nil && o.overrides.envVar != "" {
		return o.overrides.envVar
	}
	return "HCLOUD_" + strings.ReplaceAll(strings.ToUpper(o.Name), "-", "_")
}

func (o *Option[T]) FlagName() string {
	if !o.HasFlag(OptionFlagPFlag) {
		return ""
	}
	if o.overrides != nil && o.overrides.flagName != "" {
		return o.overrides.flagName
	}
	return "--" + o.Name
}

func (o *Option[T]) Completions() []string {
	var t T
	switch any(t).(type) {
	case bool:
		return []string{"true", "false"}
	}
	return nil
}

func (o *Option[T]) T() any {
	var t T
	return t
}

func (o *Option[T]) addToFlagSet(fs *pflag.FlagSet) {
	if !o.HasFlag(OptionFlagPFlag) {
		return
	}
	switch v := any(o.Default).(type) {
	case bool:
		fs.Bool(o.Name, v, o.Description)
	case string:
		fs.String(o.Name, v, o.Description)
	case time.Duration:
		fs.Duration(o.Name, v, o.Description)
	case []string:
		fs.StringSlice(o.Name, v, o.Description)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

func newOpt[T any](name, description string, def T, flags OptionFlag, ov *overrides) *Option[T] {
	o := &Option[T]{Name: name, Description: description, Default: def, Flags: flags, overrides: ov}
	Options[name] = o
	return o
}

// NewTestOption is a helper function to create an option for testing purposes
func NewTestOption[T any](name, description string, def T, flags OptionFlag, ov *overrides) (*Option[T], func()) {
	opt := newOpt(name, description, def, flags, ov)
	return opt, func() {
		delete(Options, name)
	}
}
