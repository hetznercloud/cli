package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type OptionSource int

const (
	// OptionSourcePreference indicates that the option can be set in the config file, globally or per context (in the preferences section)
	OptionSourcePreference OptionSource = 1 << iota
	// OptionSourceConfig indicates that the option can be set in the config file, but only globally or per context (not in the preferences section)
	OptionSourceConfig
	// OptionSourceFlag indicates that the option can be set via a command line flag
	OptionSourceFlag
	// OptionSourceEnv indicates that the option can be set via an environment variable
	OptionSourceEnv
)

type opt interface {
	AddToFlagSet(fs *pflag.FlagSet)
	HasSource(src OptionSource) bool
	T() any
}

var opts = make(map[string]opt)

var (
	OptionConfig         = newOpt("config", "Config file path", DefaultConfigPath(), OptionSourceFlag|OptionSourceEnv)
	OptionToken          = newOpt("token", "Hetzner Cloud API token", "", OptionSourceConfig|OptionSourceEnv)
	OptionEndpoint       = newOpt("endpoint", "Hetzner Cloud API endpoint", "", OptionSourcePreference|OptionSourceFlag|OptionSourceEnv)
	OptionDebug          = newOpt("debug", "Enable debug output", false, OptionSourcePreference|OptionSourceFlag|OptionSourceEnv)
	OptionDebugFile      = newOpt("debug-file", "Write debug output to file", "", OptionSourcePreference|OptionSourceFlag|OptionSourceEnv)
	OptionContext        = newOpt("context", "Active context", "", OptionSourceConfig|OptionSourceFlag|OptionSourceEnv)
	OptionPollInterval   = newOpt("poll-interval", "Interval at which to poll information, for example action progress", 500*time.Millisecond, OptionSourcePreference|OptionSourceFlag|OptionSourceEnv)
	OptionQuiet          = newOpt("quiet", "Only print error messages", false, OptionSourcePreference|OptionSourceFlag|OptionSourceEnv)
	OptionDefaultSSHKeys = newOpt("default-ssh-keys", "Default SSH keys for new servers", []string{}, OptionSourcePreference|OptionSourceEnv)
	OptionSSHPath        = newOpt("ssh-path", "Path to the ssh binary", "ssh", OptionSourcePreference|OptionSourceFlag|OptionSourceEnv)
)

type Option[T any] struct {
	Name    string
	Usage   string
	Default T
	Source  OptionSource
}

func (o *Option[T]) Value() T {
	return viper.Get(o.Name).(T)
}

func (o *Option[T]) SetValue(v T) {
	viper.Set(o.Name, v)
}

func (o *Option[T]) HasSource(src OptionSource) bool {
	return o.Source&src != 0
}

func (o *Option[T]) T() any {
	var t T
	return t
}

func (o *Option[T]) AddToFlagSet(fs *pflag.FlagSet) {
	if !o.HasSource(OptionSourceFlag) {
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

func newOpt[T any](name, usage string, def T, source OptionSource) *Option[T] {
	o := &Option[T]{Name: name, Usage: usage, Default: def, Source: source}
	opts[name] = o
	viper.SetDefault(name, def)
	return o
}
