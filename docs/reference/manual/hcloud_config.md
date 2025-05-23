## hcloud config

Manage configuration

### Synopsis

This command allows you to manage options for the Hetzner Cloud CLI. Options can be set inside the
configuration file, through environment variables or with flags. 

The hierarchy for configuration sources is as follows (from highest to lowest priority):
1. Flags
2. Environment variables
3. Configuration file (context)
4. Configuration file (global)
5. Default values

Option values can have following types:
 - string
 - integer
 - boolean (true/false, yes/no)
 - duration (in the Go duration format, e.g. "1h30m")
 - any of the above as a list

Most options are 'preferences' - these options can be set globally and can additionally be overridden
for each context. Below is a list of all non-preference options:

```
┌─────────┬──────────────────────┬────────┬────────────────┬──────────────────────┬───────────┐
│ OPTION  │ DESCRIPTION          │ TYPE   │ CONFIG KEY     │ ENVIRONMENT VARIABLE │ FLAG      │
├─────────┼──────────────────────┼────────┼────────────────┼──────────────────────┼───────────┤
│ config  │ Config file path     │ string │                │ HCLOUD_CONFIG        │ --config  │
│         │ (default             │        │                │                      │           │
│         │ "~/.config/hcloud/cl │        │                │                      │           │
│         │ i.toml")             │        │                │                      │           │
├─────────┼──────────────────────┼────────┼────────────────┼──────────────────────┼───────────┤
│ context │ Currently active     │ string │ active_context │ HCLOUD_CONTEXT       │ --context │
│         │ context              │        │                │                      │           │
├─────────┼──────────────────────┼────────┼────────────────┼──────────────────────┼───────────┤
│ token   │ Hetzner Cloud API    │ string │ token          │ HCLOUD_TOKEN         │           │
│         │ token                │        │                │                      │           │
└─────────┴──────────────────────┴────────┴────────────────┴──────────────────────┴───────────┘
```

Since the above options are not preferences, they cannot be modified with 'hcloud config set' or 
'hcloud config unset'. However, you are able to retrieve them using 'hcloud config get' and 'hcloud config list'.
Following options are preferences and can be used with set/unset/add/remove:

```
┌──────────────────┬──────────────────────┬─────────────┬──────────────────┬─────────────────────────┬─────────────────┐
│ OPTION           │ DESCRIPTION          │ TYPE        │ CONFIG KEY       │ ENVIRONMENT VARIABLE    │ FLAG            │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ debug            │ Enable debug output  │ boolean     │ debug            │ HCLOUD_DEBUG            │ --debug         │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ debug-file       │ File to write debug  │ string      │ debug_file       │ HCLOUD_DEBUG_FILE       │ --debug-file    │
│                  │ output to            │             │                  │                         │                 │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ default-ssh-keys │ Default SSH Keys for │ string list │ default_ssh_keys │ HCLOUD_DEFAULT_SSH_KEYS │                 │
│                  │ new Servers          │             │                  │                         │                 │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ endpoint         │ Hetzner Cloud API    │ string      │ endpoint         │ HCLOUD_ENDPOINT         │ --endpoint      │
│                  │ endpoint             │             │                  │                         │                 │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ poll-interval    │ Interval at which to │ duration    │ poll_interval    │ HCLOUD_POLL_INTERVAL    │ --poll-interval │
│                  │ poll information,    │             │                  │                         │                 │
│                  │ for example action   │             │                  │                         │                 │
│                  │ progress             │             │                  │                         │                 │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ quiet            │ If true, only print  │ boolean     │ quiet            │ HCLOUD_QUIET            │ --quiet         │
│                  │ error messages       │             │                  │                         │                 │
├──────────────────┼──────────────────────┼─────────────┼──────────────────┼─────────────────────────┼─────────────────┤
│ sort.<resource>  │ Default sorting for  │ string list │ sort.<resource>  │ HCLOUD_SORT_<RESOURCE>  │                 │
│                  │ resource             │             │                  │                         │                 │
└──────────────────┴──────────────────────┴─────────────┴──────────────────┴─────────────────────────┴─────────────────┘
```

Options will be persisted in the configuration file. To find out where your configuration file is located
on disk, run 'hcloud config get config'.


### Options

```
  -h, --help   help for config
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
* [hcloud config add](hcloud_config_add.md)	 - Add values to a list
* [hcloud config get](hcloud_config_get.md)	 - Get a configuration value
* [hcloud config list](hcloud_config_list.md)	 - List configuration values
* [hcloud config remove](hcloud_config_remove.md)	 - Remove values from a list
* [hcloud config set](hcloud_config_set.md)	 - Set a configuration value
* [hcloud config unset](hcloud_config_unset.md)	 - Unset a configuration value
