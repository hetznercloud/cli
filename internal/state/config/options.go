package config

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type OptionFlag int

const (
	// OptionFlagPreference indicates that the option can be set in the config file, globally or per context (in the preferences section)
	OptionFlagPreference OptionFlag = 1 << iota
	// OptionFlagConfig indicates that the option can be set in the config file, but only globally or per context (not in the preferences section)
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
	HasFlag(src OptionFlag) bool
	GetAsAny(c Config) any
	OverrideAny(c Config, v any)
	Changed(c Config) bool
	Completions() []string
	IsSlice() bool
	T() any
}

var Options = make(map[string]IOption)

// Note: &^ is the bit clear operator and is used to remove flags from the default flag set
var (
	OptionConfig         = newOpt("config", "Config file path", DefaultConfigPath(), OptionFlagPFlag|OptionFlagEnv)
	OptionToken          = newOpt("token", "Hetzner Cloud API token", "", OptionFlagConfig|OptionFlagEnv|OptionFlagSensitive)
	OptionContext        = newOpt("context", "Active context", "", OptionFlagConfig|OptionFlagEnv|OptionFlagPFlag)
	OptionEndpoint       = newOpt("endpoint", "Hetzner Cloud API endpoint", hcloud.Endpoint, DefaultPreferenceFlags)
	OptionDebug          = newOpt("debug", "Enable debug output", false, DefaultPreferenceFlags)
	OptionDebugFile      = newOpt("debug-file", "Write debug output to file", "", DefaultPreferenceFlags)
	OptionPollInterval   = newOpt("poll-interval", "Interval at which to poll information, for example action progress", 500*time.Millisecond, DefaultPreferenceFlags)
	OptionQuiet          = newOpt("quiet", "Only print error messages", false, DefaultPreferenceFlags)
	OptionDefaultSSHKeys = newOpt("default-ssh-keys", "Default SSH keys for new servers", []string{}, DefaultPreferenceFlags&^OptionFlagPFlag)
	OptionSSHPath        = newOpt("ssh-path", "Path to the ssh binary", "ssh", DefaultPreferenceFlags)
)

type Option[T any] struct {
	Name    string
	Usage   string
	Default T
	Source  OptionFlag
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
			val = cast.ToStringSlice(v)
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
	return o.Source&src != 0
}

func (o *Option[T]) IsSlice() bool {
	return reflect.TypeOf(o.T()).Kind() == reflect.Slice
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
		fs.Bool(o.Name, v, o.Usage)
	case string:
		fs.String(o.Name, v, o.Usage)
	case time.Duration:
		fs.Duration(o.Name, v, o.Usage)
	case []string:
		fs.StringSlice(o.Name, v, o.Usage)
	default:
		panic(fmt.Sprintf("unsupported type %T", v))
	}
}

func newOpt[T any](name, usage string, def T, source OptionFlag) *Option[T] {
	o := &Option[T]{Name: name, Usage: usage, Default: def, Source: source}
	Options[name] = o
	return o
}
