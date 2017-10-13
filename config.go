package cli

import toml "github.com/pelletier/go-toml"

type Config struct {
	Token    string
	Endpoint string
}

func UnmarshalConfig(data []byte) (*Config, error) {
	var v struct {
		CLI struct {
			Token    string `toml:"token"`
			Endpoint string `toml:"endpoint"`
		} `toml:"cli"`
	}
	if err := toml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	config := &Config{
		Token:    v.CLI.Token,
		Endpoint: v.CLI.Endpoint,
	}
	return config, nil
}
