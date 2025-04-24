## hcloud completion

Output shell completion code for the specified shell

### Synopsis

To load completions:

### Bash

To load completions into the current shell execute:

    source <(hcloud completion bash)

In order to make the completions permanent, append the line above to
your .bashrc.

### Zsh

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


```
hcloud completion <shell>
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string            Config file path (default "~/.config/hcloud/cli.toml")
      --context string           Currently active context
      --debug                    Enable debug output
      --debug-file string        File to write debug output to
      --endpoint string          Hetzner Cloud API endpoint (default "https://api.hetzner.cloud/v1")
      --poll-interval duration   Interval at which to poll information, for example action progress (default 500ms)
      --quiet                    If true, only print error messages
```

### SEE ALSO

* [hcloud](hcloud.md)	 - Hetzner Cloud CLI
