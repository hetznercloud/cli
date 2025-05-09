# Setup the hcloud CLI

This tutorial will guide you through the process of setting up the hcloud CLI on your local machine.

## Prerequisites

Before you begin, ensure you have the following:

- A [Hetzner Cloud account](https://console.hetzner.cloud).

## 1. Install the hcloud CLI

### 1.1 Manual installation

You can download pre-built binaries from our [GitHub releases](https://github.com/hetznercloud/cli/releases). 
Install them by extracting the archive and moving the binary to a directory in your `PATH`.

On a 64-bit Linux system, it could look something like this:

```bash
curl -LO https://github.com/hetznercloud/cli/releases/download/v1.50.0/hcloud-darwin-amd64.tar.gz
sudo tar -C /usr/local/bin -xzf hcloud-darwin-amd64.tar.gz
rm hcloud-darwin-amd64.tar.gz
```

### 1.2 Installation using Go

If you have Go installed, you can also install hcloud CLI from source using the following command:

```bash
go install github.com/hetznercloud/cli@latest
```

> [!NOTE]
> Binaries built with Go will not have the correct version embedded.

> [!NOTE]
> Both of the above installation methods do not provide automatic updates. Please make sure to keep your installation up to date manually.

### 1.3 Installation using Homebrew

On Linux and macOS you can also install the hcloud CLI using Homebrew:

```bash
brew install hcloud
```

### 1.4 Installation using scoop

On Windows, you can install `hcloud` using scoop:

```bash
scoop install hcloud
```

### 1.5 Using hcloud with Docker

Instead of installing hcloud on the host, you can also use our docker image at `hetznercloud/cli`.

```bash
docker run --rm -e HCLOUD_TOKEN="<your token>" hetznercloud/cli:latest <command>
```

If you want to use (and persist) your configuration, you can mount it to `/config.toml`:
```bash
docker run --rm -v ~/.config/hcloud/cli.toml:/config.toml hetznercloud/cli:latest <command>
```

The image is based on Alpine Linux, so a shell is available in the image. You can use it to run commands interactively:

```bash
docker run -it --rm --entrypoint /bin/sh hetznercloud/cli:latest
```

---

> [!WARNING]
> Debian-based distributions (using apt) provide outdated versions of the hcloud CLI.
> Please consider one of the other installation methods.

## 2. (Optional) Setup auto-completion

hcloud CLI offers auto-completion for bash, zsh, fish and PowerShell. It is recommended to enable it for a better user experience.

### 2.1 Bash

To load completions into the current shell execute:

    source <(hcloud completion bash)

In order to make the completions permanent, append the line above to
your .bashrc.

### 2.2 Zsh

If shell completions are not already enabled for your environment need
to enable them. Add the following line to your ~/.zshrc file:

    autoload -Uz compinit; compinit

To load completions for each session execute the following commands:

    mkdir -p ~/.config/hcloud/completion/zsh
    hcloud completion zsh > ~/.config/hcloud/completion/zsh/_hcloud

Finally add the following line to your ~/.zshrc file, *before* you
call the compinit function:

    fpath+=(~/.config/hcloud/completion/zsh)

In the end your ~/.zshrc file should contain the following two lines
in the order given here.

    fpath+=(~/.config/hcloud/completion/zsh)
    #  ... anything else that needs to be done before compinit
    autoload -Uz compinit; compinit
    # ...

You will need to start a new shell for this setup to take effect.

### 2.3 Fish

To load completions into the current shell execute:

    hcloud completion fish | source

In order to make the completions permanent execute once:

    hcloud completion fish > ~/.config/fish/completions/hcloud.fish

### 2.4 PowerShell

To load completions into the current shell execute:

    PS> hcloud completion powershell | Out-String | Invoke-Expression

To load completions for every new session, run
and source this file from your PowerShell profile.

    PS> hcloud completion powershell > hcloud.ps1

## 3. Create a context

The hcloud CLI uses contexts to manage multiple Hetzner Cloud tokens and set configuration preferences.

First, you need to create an API token.
Follow the instructions in the [Hetzner Cloud documentation](https://docs.hetzner.com/cloud/api/getting-started/generating-api-token) to create your project API token.

Once you have your token, you can create a context using the following command:

```bash
hcloud context create <context-name>
```

Ideally, keep the context name similar to the project name so you can easily identify it later.

When prompted, enter the API token you created earlier. Your context should now be created and activated.

## 4. Verify everything works

To verify that the hcloud CLI is working correctly, you can run the following command:

```bash
hcloud datacenter list
```

You should see something like this:

```plaintext
ID   NAME        DESCRIPTION                   LOCATION
2    nbg1-dc3    Nuremberg 1 virtual DC 3      nbg1    
3    hel1-dc2    Helsinki 1 virtual DC 2       hel1    
4    fsn1-dc14   Falkenstein 1 virtual DC 14   fsn1    
5    ash-dc1     Ashburn virtual DC 1          ash     
6    hil-dc1     Hillsboro virtual DC 1        hil     
7    sin-dc1     Singapore virtual DC 1        sin     
```

If you see this output, congratulations! You have successfully set up the hcloud CLI on your local machine.

If there are any problems, make sure you followed all steps of this tutorial correctly. If there still are problems,
you can reach out to our [Support](https://console.hetzner.cloud/support) to get help.
