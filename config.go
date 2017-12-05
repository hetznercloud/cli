package cli

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml"
)

var DefaultConfigPath string

func init() {
	if home := os.Getenv("HOME"); home != "" {
		DefaultConfigPath = filepath.Join(home, ".config", "hcloud", "config.toml")
	}
}

type Config struct {
	Token    string
	Endpoint string
}

type RawConfig struct {
	CLI struct {
		Token    string `toml:"token,omitempty"`
		Endpoint string `toml:"endpoint,omitempty"`
	} `toml:"cli"`
}

func MarshalConfig(c *Config) ([]byte, error) {
	if c == nil {
		return []byte{}, nil
	}

	var raw RawConfig
	raw.CLI.Token = c.Token
	raw.CLI.Endpoint = c.Endpoint
	return toml.Marshal(raw)
}

func UnmarshalConfig(data []byte) (*Config, error) {
	var raw RawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &Config{
		Token:    raw.CLI.Token,
		Endpoint: raw.CLI.Endpoint,
	}, nil
}
