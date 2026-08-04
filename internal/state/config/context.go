package config

import "fmt"

type Context interface {
	Name() string
	Token() string
	Preferences() Preferences
}

func NewContext(name, token string) Context {
	return &context{
		ContextName:  name,
		ContextToken: token,
	}
}

type context struct {
	ContextName        string      `toml:"name"`
	ContextToken       string      `toml:"token"`
	ContextPreferences Preferences `toml:"preferences"`
}

func (ctx *context) Name() string {
	return ctx.ContextName
}

// Token returns the token for the context.
// If you just need the token regardless of the context, please use [OptionToken] instead.
func (ctx *context) Token() string {
	return ctx.ContextToken
}

func (ctx *context) Preferences() Preferences {
	if ctx.ContextPreferences == nil {
		ctx.ContextPreferences = make(Preferences)
	}
	return ctx.ContextPreferences
}

func ContextNames(cfg Config) []string {
	ctxs := cfg.Contexts()
	names := make([]string, len(ctxs))
	for i, ctx := range ctxs {
		names[i] = ctx.Name()
	}
	return names
}

func ContextByName(cfg Config, name string) Context {
	for _, c := range cfg.Contexts() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func RemoveContext(cfg Config, context Context) error {
	var filtered []Context
	for _, c := range cfg.Contexts() {
		if c != context {
			filtered = append(filtered, c)
		}
	}
	return cfg.SetContexts(filtered)
}

func RenameContext(ctx Context, newName string) error {
	internalContext, ok := ctx.(*context)
	if !ok {
		return fmt.Errorf("unsupported context implementation %T", ctx)
	}
	internalContext.ContextName = newName
	return nil
}
