package config

type Context interface {
	Name() string
	Token() string
	TokenCommand() string
	Preferences() Preferences
}

func NewContext(name, token string) Context {
	return &context{
		ContextName:  name,
		ContextToken: token,
	}
}

// NewTokenCommandContext creates a context that retrieves its token from a command.
func NewTokenCommandContext(name, command string) Context {
	return &context{ContextName: name, ContextTokenCommand: command}
}

type context struct {
	ContextName         string      `toml:"name"`
	ContextToken        string      `toml:"token,omitempty"`
	ContextTokenCommand string      `toml:"token_command,omitempty"`
	ContextPreferences  Preferences `toml:"preferences"`
}

func (ctx *context) Name() string {
	return ctx.ContextName
}

// Token returns the token for the context.
// If you just need the token regardless of the context, please use [OptionToken] instead.
func (ctx *context) Token() string {
	return ctx.ContextToken
}

// TokenCommand returns the shell command used to retrieve this context's token.
func (ctx *context) TokenCommand() string {
	return ctx.ContextTokenCommand
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

func RemoveContext(cfg Config, context Context) {
	var filtered []Context
	for _, c := range cfg.Contexts() {
		if c != context {
			filtered = append(filtered, c)
		}
	}
	cfg.SetContexts(filtered)
}

func RenameContext(ctx Context, newName string) {
	if ctx, ok := (ctx).(*context); ok {
		ctx.ContextName = newName
	}
}
