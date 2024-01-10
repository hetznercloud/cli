package context_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/cmd/context"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestSSHKeyAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(nil)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	activeContext := &config.Context{
		Name:    "test",
		SSHKeys: []string{"existing_key"},
	}

	fx.Config.EXPECT().
		ActiveContext().
		Return(activeContext)
	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "new_key").
		Return(&hcloud.SSHKey{}, nil, nil)
	fx.Config.EXPECT().
		Write()

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "add", "new_key"})

	assert.NoError(t, err)
	assert.Equal(t, "Added SSH key(s) new_key to context \"test\"\n", out)
	assert.Empty(t, errOut)

	assert.Equal(t, []string{"existing_key", "new_key"}, activeContext.SSHKeys)
}

func TestSSHKeyAddAll(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(nil)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	activeContext := &config.Context{
		Name:    "test",
		SSHKeys: []string{"existing_key"},
	}

	fx.Config.EXPECT().
		ActiveContext().
		Return(activeContext)
	fx.Client.SSHKeyClient.EXPECT().
		All(gomock.Any()).
		Return([]*hcloud.SSHKey{{ID: 42, Name: "foo"}, {ID: 1337, Name: "bar"}}, nil)
	fx.Config.EXPECT().
		Write()

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "add", "--all"})

	assert.NoError(t, err)
	assert.Equal(t, "Added SSH key(s) 42, 1337 to context \"test\"\n", out)
	assert.Empty(t, errOut)

	assert.Equal(t, []string{"1337", "42", "existing_key"}, activeContext.SSHKeys)
}

func TestSSHKeyAddContext(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	contexts := []*config.Context{
		{
			Name: "test",
		},
		{
			Name:    "test2",
			SSHKeys: []string{"existing_key"},
		},
	}

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(contexts)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	fx.Client.SSHKeyClient.EXPECT().
		Get(gomock.Any(), "new_key").
		Return(&hcloud.SSHKey{}, nil, nil)
	fx.Config.EXPECT().
		Write()

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "add", "new_key", "--context", "test2"})

	assert.NoError(t, err)
	assert.Equal(t, "Added SSH key(s) new_key to context \"test2\"\n", out)
	assert.Empty(t, errOut)

	assert.Equal(t, []string{"existing_key", "new_key"}, contexts[1].SSHKeys)
}

func TestSSHKeyRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(nil)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	activeContext := &config.Context{
		Name:    "test",
		SSHKeys: []string{"remove_me", "dont_remove_me"},
	}

	fx.Config.EXPECT().
		ActiveContext().
		Return(activeContext)
	fx.Config.EXPECT().
		Write()

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "remove", "remove_me"})

	assert.NoError(t, err)
	assert.Equal(t, "Removed 1 SSH key(s) from context \"test\"\n", out)
	assert.Empty(t, errOut)

	assert.Equal(t, []string{"dont_remove_me"}, activeContext.SSHKeys)
}

func TestSSHKeyRemoveContext(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	contexts := []*config.Context{
		{
			Name: "test",
		},
		{
			Name:    "test2",
			SSHKeys: []string{"remove_me", "dont_remove_me"},
		},
	}

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(contexts)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	fx.Config.EXPECT().
		Write()

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "remove", "remove_me", "--context", "test2"})

	assert.NoError(t, err)
	assert.Equal(t, "Removed 1 SSH key(s) from context \"test2\"\n", out)
	assert.Empty(t, errOut)

	assert.Equal(t, []string{"dont_remove_me"}, contexts[1].SSHKeys)
}

func TestSSHKeyRemoveAll(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(nil)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	activeContext := &config.Context{
		Name:    "test",
		SSHKeys: []string{"remove_me", "remove_me_too"},
	}

	fx.Config.EXPECT().
		ActiveContext().
		Return(activeContext)
	fx.Config.EXPECT().
		Write()

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "remove", "--all"})

	assert.NoError(t, err)
	assert.Equal(t, "Removed 2 SSH key(s) from context \"test\"\n", out)
	assert.Empty(t, errOut)

	assert.Empty(t, activeContext.SSHKeys)
}

func TestSSHKeyList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	rootCmd := cli.NewRootCommand(fx.State())

	// needed because subcommands fetch a list of context names for completions
	fx.Config.EXPECT().
		Contexts().
		AnyTimes().
		Return(nil)

	rootCmd.AddCommand(context.NewCommand(fx.State()))

	activeContext := &config.Context{
		Name:    "test",
		SSHKeys: []string{"foo", "bar", "baz"},
	}

	fx.Config.EXPECT().
		ActiveContext().
		Return(activeContext)

	out, errOut, err := fx.Run(rootCmd, []string{"context", "ssh-key", "list"})

	assert.NoError(t, err)
	assert.Equal(t, "SSH keys in context \"test\":\n - foo\n - bar\n - baz\n", out)
	assert.Empty(t, errOut)
}
