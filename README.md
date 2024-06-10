# hcloud: Command-line interface for Hetzner Cloud

[![Release](https://img.shields.io/github/v/release/hetznercloud/cli)](https://github.com/hetznercloud/cli/releases/latest)
![Go Version](https://img.shields.io/github/go-mod/go-version/hetznercloud/cli/main?label=Go)
[![CI](https://github.com/hetznercloud/cli/actions/workflows/ci.yml/badge.svg)](https://github.com/hetznercloud/cli/actions/workflows/ci.yml)
[![Build](https://github.com/hetznercloud/cli/actions/workflows/build.yml/badge.svg)](https://github.com/hetznercloud/cli/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/hetznercloud/cli/graph/badge.svg?token=fFDgg6Ua6U)](https://codecov.io/gh/hetznercloud/cli)

`hcloud` is a command-line interface for interacting with Hetzner Cloud.

[![asciicast](https://asciinema.org/a/157991.png)](https://asciinema.org/a/157991)

## Installation

You can download pre-built binaries for Linux, FreeBSD, macOS, and Windows on
the [releases page](https://github.com/hetznercloud/cli/releases).

On macOS and Linux, you can install `hcloud` via [Homebrew](https://brew.sh/):

    brew install hcloud


On Windows, you can install `hcloud` via [Scoop](https://scoop.sh/)

    scoop install hcloud

### Third-party packages

There are unofficial packages maintained by third-party users. Please note
that these packages aren’t supported nor maintained by Hetzner Cloud and
may not always be up-to-date. Downloading the binary or building from source
is still the recommended install method.

### Build manually

If you have Go installed, you can build and install the latest version of
`hcloud` with:

    go install github.com/hetznercloud/cli/cmd/hcloud@latest

> Binaries built in this way do not have the correct version embedded. Use our
prebuilt binaries or check out [`.goreleaser.yml`](.goreleaser.yml) to learn
how to embed it yourself.

## Getting Started

1.  Visit the Hetzner Cloud Console at [console.hetzner.cloud](https://console.hetzner.cloud/),
    select your project, and create a new API token.

2.  Configure the `hcloud` program to use your token:

        hcloud context create my-project

3.  You’re ready to use the program. For example, to get a list of available server
    types, run:

        hcloud server-type list

See `hcloud help` for a list of commands.

## Shell Completion

`hcloud` provides completions for various shells.

### Bash

To load completions into the current shell execute:

    source <(hcloud completion bash)

In order to make the completions permanent, append the line above to
your `.bashrc`.

### Zsh

If shell completions are not already enabled for your environment need
to enable them. Add the following line to your `~/.zshrc` file:

    autoload -Uz compinit; compinit

To load completions for each session execute the following commands:

    mkdir -p ~/.config/hcloud/completion/zsh
    hcloud completion zsh > ~/.config/hcloud/completion/zsh/_hcloud

Finally, add the following line to your `~/.zshrc` file, *before* you
call the `compinit` function:

    fpath+=(~/.config/hcloud/completion/zsh)

In the end your `~/.zshrc` file should contain the following two lines
in the order given here.

    fpath+=(~/.config/hcloud/completion/zsh)
    #  ... anything else that needs to be done before compinit
    autoload -Uz compinit; compinit
    # ...

You will need to start a new shell for this setup to take effect.

### Fish

To load completions into the current shell execute:

    hcloud completion fish | source

In order to make the completions permanent execute once:

     hcloud completion fish > ~/.config/fish/completions/hcloud.fish

### PowerShell:

To load completions into the current shell execute:

    PS> hcloud completion powershell | Out-String | Invoke-Expression

To load completions for every new session, run
and source this file from your PowerShell profile.

    PS> hcloud completion powershell > hcloud.ps1

## Output configuration

You can control output via the `-o` option:

* For `list` commands, you can specify `-o noheader` to omit the table header.

* For `list` commands, you can specify `-o columns=id,name` to only show certain
  columns in the table.

* For `describe` commands, you can specify `-o json` to get a JSON representation
  of the resource. The schema is identical to those in the Hetzner Cloud API which
  are documented at [docs.hetzner.cloud](https://docs.hetzner.cloud).

* For `create` commands, you can specify `-o json` to get a JSON representation
  of the API response. API responses are documented at [docs.hetzner.cloud](https://docs.hetzner.cloud).
  In contrast to `describe` commands, `create` commands can return extra information, for example
  the initial root password of a server.

* For `describe` commands, you can specify `-o format={{.ID}}` to format output
  according to the given [Go template](https://golang.org/pkg/text/template/).
  The template’s input is the resource’s corresponding struct in the
  [hcloud-go](https://godoc.org/github.com/hetznercloud/hcloud-go/hcloud) library.

## Configuring hcloud

The hcloud CLI tool can be configured using following methods:
1. Configuration file
2. Environment variables
3. Command line flags

A higher number means a higher priority. For example, a command line flag will 
always override an environment variable.

The configuration file is located at `~/.config/hcloud/cli.toml` by default 
(On Windows: `%APPDATA%\hcloud\cli.toml`). You can change the location by setting 
the `HCLOUD_CONFIG` environment variable or the `--config` flag. The configuration file
stores global preferences, the currently active context, all contexts and 
context-specific preferences. Contexts always store a token and can optionally have
additional preferences which take precedence over the globally set preferences.

However, a config file is not required. If no config file is found, the CLI will
use the default configuration. Overriding options using environment variables allows the
hcloud CLI to function in a stateless way. For example, setting `HCLOUD_TOKEN` is
already enough in many cases.

You can use the `hcloud config` command to manage your configuration, for example 
to get, list, set and unset configuration options and preferences. You can view a list
of all available options and preferences by running `hcloud config --help`.

## Examples

### List all servers

```
$ hcloud server list
ID       NAME                    STATUS    IPV4
210216   test1                   running   78.46.122.12
210729   ubuntu-8gb-nbg1-dc3-1   running   94.130.177.158
```

### Create a server

```
$ hcloud server create --name test --image debian-9 --type cx22 --ssh-key demo
   7s [====================================================================] 100%
Server 325211 created
```

## License

MIT license
