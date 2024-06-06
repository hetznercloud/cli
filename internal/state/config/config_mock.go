package config

import (
	"io"
	"os"
)

// We do not need to generate a gomock for the Config, since you can set config
// values during tests with viper.Set()

type MockConfig struct {
	Config
}

func (c *MockConfig) Write(_ io.Writer) error {
	// MockConfig always writes to stdout for testing purposes
	return c.Config.Write(os.Stdout)
}

var _ Config = (*MockConfig)(nil)
