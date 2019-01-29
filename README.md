# hcloud: Command-line interface for Hetzner Cloud

[![Build status](https://travis-ci.org/hetznercloud/cli.svg?branch=master)](https://travis-ci.org/hetznercloud/cli)

`hcloud` is a command-line interface for interacting with Hetzner Cloud.

[![asciicast](https://asciinema.org/a/157991.png)](https://asciinema.org/a/157991)

## Installation

You can download pre-built binaries for Linux, FreeBSD, macOS, and Windows on
the [releases page](https://github.com/hetznercloud/cli/releases).

On macOS, you can install `hcloud` via [Homebrew](https://brew.sh/):

    brew install hcloud

### Third-party packages

There are inofficial packages maintained by third-party users. Please note
that these packages aren’t supported nor maintained by Hetzner Cloud and
may not always be up-to-date. Downloading the binary or building from source
is still the recommended install method.

| Operating System | Command                                           |
| ---------------- | ------------------------------------------------- |
| Arch Linux       | `pacman -Syu hcloud`                              |
| Void Linux       | `xbps-install -Syu hcloud`                        |
 
### Build manually

If you have Go installed, you can build and install the `hcloud` program with:

    go get -u github.com/hetznercloud/cli/cmd/hcloud

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

To enable shell completion, run one of the following commands (or better,
add it to your `.bashrc` or `.zshrc`):

```
$ source <(hcloud completion bash)   # bash
$ source <(hcloud completion zsh)    # zsh
```

## Configure hcloud using environment variables

You can use the following environment variables to configure `hcloud`:

* `HCLOUD_TOKEN`
* `HCLOUD_CONTEXT`
* `HCLOUD_CONFIG`

When using `hcloud` in scripts, for example, it may be cumbersome to work with
contexts. Instead of creating a context, you can set the token via the `HCLOUD_TOKEN`
environment variable. When combined with tools like [direnv](https://direnv.net), you
can configure a per-directory context by setting `HCLOUD_CONTEXT=my-context` via `.envrc`.

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
$ hcloud server create --name test --image debian-9 --type cx11 --ssh-key demo
   7s [====================================================================] 100%
Server 325211 created
```

## License

MIT license
