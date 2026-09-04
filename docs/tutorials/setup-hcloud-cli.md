<!--
---
date: "2026-09-07"
date_changed: "2026-09-07"
title: "Installing the hcloud CLI"
tags: []
language: "en"
description: ""
docs_type: ["getting_started"]
product_category: ["Integrations"]
translation: ["Integrations", "CLI", "Getting Started", "Installing the hcloud CLI"]
scrape_type: "whole"
priority: 100
---
-->

# Setup the hcloud CLI

This tutorial will guide you through the process of setting up the hcloud CLI on your local machine.

**Prerequisites:**

Before you begin, ensure you have the following:

- A [Hetzner Console account](https://console.hetzner.com).

<br>

----

<br>

1. Install the hcloud CLI
   
   Click on one of the options below to view the respective steps.
   
   <details>
   <summary>Manual installation</summary>
   <ul>
   
   You can download pre-built binaries from our [GitHub releases](https://github.com/hetznercloud/cli/releases). 
   Install them by extracting the archive and moving the binary to a directory in your `PATH`.
   
   On a 64-bit Linux system, it could look something like this:
   
   ```bash
   curl -sSLO https://github.com/hetznercloud/cli/releases/latest/download/hcloud-linux-amd64.tar.gz
   sudo tar -C /usr/local/bin --no-same-owner -xzf hcloud-linux-amd64.tar.gz hcloud
   rm hcloud-linux-amd64.tar.gz
   ```
   
   <br></ul>
   </details>
   <details>
   <summary>Installation using Go</summary>
   <ul>
   
   If you have Go installed, you can also install hcloud CLI from source using the following command:
   
   ```bash
   go install github.com/hetznercloud/cli/cmd/hcloud@latest
   ```
   > [!NOTE]
   > Binaries built with Go will not have the correct version embedded.
   
   > [!NOTE]
   > Both of the above installation methods do not provide automatic updates. Please make sure to keep your installation up to date manually.
   
   <br></ul>
   </details>
   <details>
   <summary>Installation using .deb package (Debian/Ubuntu)</summary>
   <ul>
   
   > [!WARNING]
   > Debian packages are experimental and might change or break in the future.
   
   On Debian and Ubuntu-based distributions, you can install the hcloud CLI using the `.deb` package:
   
   1. Download the latest `.deb` package from the [GitHub releases page](https://github.com/hetznercloud/cli/releases/latest).
   
   2. Install the package using `dpkg`:
   
   ```bash
   # Download the package (replace with the actual version and architecture)
   curl -sSLO https://github.com/hetznercloud/cli/releases/latest/download/hcloud-cli_<version>_linux_amd64.deb
   
   # Install the package
   sudo dpkg -i hcloud-cli_<version>_linux_amd64.deb
   ```
   
   The `.deb` package includes shell completions for bash, zsh, and fish, as well as man pages.
   
   <br></ul>
   </details>
   <details>
   <summary>Installation using .rpm package (Fedora/RHEL/CentOS)</summary>
   <ul>
   
   > [!WARNING]
   > RPM packages are experimental and might change or break in the future.
   
   On Fedora, RHEL, CentOS, and other RPM-based distributions, you can install the hcloud CLI using the `.rpm` package:
   
   1. Download the latest `.rpm` package from the [GitHub releases page](https://github.com/hetznercloud/cli/releases/latest).
   
   2. Install the package using `rpm` or `dnf`:
   
   ```bash
   # Download the package (replace with the actual version and architecture)
   curl -sSLO https://github.com/hetznercloud/cli/releases/latest/download/hcloud-cli-<version>-1.x86_64.rpm
   
   # Install using dnf (Fedora/RHEL 8+/CentOS 8+)
   sudo dnf install hcloud-cli-<version>-1.x86_64.rpm
   
   # Or install using rpm
   sudo rpm -i hcloud-cli-<version>-1.x86_64.rpm
   ```
   
   The `.rpm` package includes shell completions for bash, zsh, and fish, as well as man pages.
   
   <br></ul>
   </details>
   <details>
   <summary>Installation using Homebrew</summary>
   <ul>
   
   On Linux and macOS you can also install the hcloud CLI using Homebrew:
   
   ```bash
   brew install hcloud
   ```
   
   <br></ul>
   </details>
   <details>
   <summary>Installation on Windows using WinGet or Scoop</summary>
   <ul>
   
   > [!WARNING]
   > The WinGet and Scoop package entries are not maintained by Hetzner.
   
   On Windows, you can install `hcloud` using WinGet (Windows Package Manager) or Scoop:
   
   ```bash
   winget install HetznerCloud.CLI
   ```
   
   ```bash
   scoop install hcloud
   ```
   
   <br></ul>
   </details>
   <details>
   <summary>Using hcloud with Docker</summary>
   <ul>
   
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
   > Some third-party package repositories may provide outdated versions of the hcloud CLI.
   > Please consider one of the other installation methods.

   <br></ul>
   </details>

<br>

2. Setup auto-completion (Optional)
   
   hcloud CLI offers auto-completion for bash, zsh, fish and PowerShell. It is recommended to enable it for a better user experience.

   Click on one of the options below to view the respective steps.
   
   <details>
   <summary>Bash</summary>
   <ul>
   
   To load completions into the current shell execute:
   ```bash
   source <(hcloud completion bash)
   ```
   
   In order to make the completions permanent, append the line above to
   your `.bashrc`.
   
   <br></ul>
   </details>
   <details>
   <summary>Zsh</summary>
   <ul>
   
   If shell completions are not already enabled for your environment need
   to enable them. Add the following line to your ~/.zshrc file:
   
   ```bash
   autoload -Uz compinit; compinit
   ```
   
   To load completions for each session execute the following commands:
   
   ```bash
   mkdir -p ~/.config/hcloud/completion/zsh
   hcloud completion zsh > ~/.config/hcloud/completion/zsh/_hcloud
   ```
   
   Finally add the following line to your ~/.zshrc file, *before* you
   call the compinit function:
   
   ```bash
   fpath+=(~/.config/hcloud/completion/zsh)
   ```
   
   In the end your ~/.zshrc file should contain the following two lines
   in the order given here.
   
   ```bash
   fpath+=(~/.config/hcloud/completion/zsh)
   #  ... anything else that needs to be done before compinit
   autoload -Uz compinit; compinit
   # ...
   ```
   
   You will need to start a new shell for this setup to take effect.
   
   <br></ul>
   </details>
   <details>
   <summary>Fish</summary>
   <ul>
   
   To load completions into the current shell execute:
   
   ```bash
   hcloud completion fish | source
   ```
   
   In order to make the completions permanent execute once:
   
   ```bash
   hcloud completion fish > ~/.config/fish/completions/hcloud.fish
   ```
   
   <br></ul>
   </details>
   <details>
   <summary>PowerShell</summary>
   <ul>
   
   To load completions into the current shell execute:
   
   ```bash
   PS> hcloud completion powershell | Out-String | Invoke-Expression
   ```
   
   To load completions for every new session, run
   and source this file from your PowerShell profile.
   
   ```bash
   PS> hcloud completion powershell > hcloud.ps1
   ```
   <br></ul>
   </details>

<br>

3. Create a context
   
   The hcloud CLI uses contexts to manage multiple Hetzner Cloud tokens and set configuration preferences.
   
   First, you need to create an API token.
   Follow the instructions in the [Hetzner Cloud documentation](https://docs.hetzner.com/cloud/api/getting-started/generating-api-token/) to create your project API token.
   
   Once you have your token, you can create a context using the following command:
   
   ```bash
   hcloud context create <context-name>
   ```
   
   Ideally, keep the context name similar to the project name so you can easily identify it later.
   
   When prompted, enter the API token you created earlier. Your context should now be created and activated.

<br>

4. Verify everything works
   
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
you can reach out to our [Support](https://console.hetzner.com/support) to get help.
