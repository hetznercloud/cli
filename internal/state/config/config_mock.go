package config

import "io"

// We do not need to generate a gomock for the Config, since you can set config
// values during tests with viper.Set()

type MockConfig struct {
	activeContext Context
	contexts      []Context
}

func (*MockConfig) Write(io.Writer) error {
	return nil
}

func (*MockConfig) ParseConfig() error {
	return nil
}

func (m *MockConfig) ActiveContext() Context {
	return m.activeContext
}

func (m *MockConfig) SetActiveContext(ctx Context) {
	m.activeContext = ctx
}

func (m *MockConfig) Contexts() []Context {
	return m.contexts
}

func (m *MockConfig) SetContexts(ctxs []Context) {
	m.contexts = ctxs
}

func (*MockConfig) Preferences() Preferences {
	return preferences{}
}

var _ Config = &MockConfig{}
