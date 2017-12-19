package cli

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml"
)

var DefaultConfigPath string

func init() {
	if home := os.Getenv("HOME"); home != "" {
		DefaultConfigPath = filepath.Join(home, ".config", "hcloud", "cli.toml")
	}
}

type Config struct {
	Token    string
	Endpoint string
}

type RawConfig struct {
	Token    string `toml:"token,omitempty"`
	Endpoint string `toml:"endpoint,omitempty"`
}

func MarshalConfig(c *Config) ([]byte, error) {
	if c == nil {
		return []byte{}, nil
	}

	var raw RawConfig
	raw.Token = c.Token
	raw.Endpoint = c.Endpoint
	return toml.Marshal(raw)
}

func UnmarshalConfig(data []byte) (*Config, error) {
	var raw RawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &Config{
		Token:    raw.Token,
		Endpoint: raw.Endpoint,
	}, nil
}
