<!--
---
date: "2026-09-07"
date_changed: "2026-09-07"
title: "Configuring hcloud"
tags: []
language: "en"
description: ""
docs_type: ["getting_started"]
product_category: ["Integrations"]
translation: ["Integrations", "CLI", "Reference", "Configuring hcloud"]
scrape_type: "whole"
priority: 90
---
-->

# Configuring hcloud

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
context-specific preferences. Contexts store a token or a token command and can optionally have
additional preferences which take precedence over the globally set preferences.

However, a config file is not required. If no config file is found, the CLI will
use the default configuration. Overriding options using environment variables allows the
hcloud CLI to function in a stateless way. For example, setting `HCLOUD_TOKEN` is
already enough in many cases.

You can use the `hcloud config` command to manage your configuration, for example
to get, list, set and unset configuration options and preferences. You can view a list
of all available options and preferences [here](manual/hcloud_config.md#synopsis).

## Token commands

A context can retrieve its API token from an external command instead of storing
it in `cli.toml`. The command must print only the token to stdout. Leading and
trailing whitespace, including newlines and carriage returns, is removed.
Diagnostics and interactive prompts must use stderr or `/dev/tty`.

```sh
hcloud context create \
  --token-command 'your-secret-manager print-token' \
  production
```

This creates the following context in `~/.config/hcloud/cli.toml`:

```toml
active_context = "production"

[[contexts]]
name = "production"
token_command = "your-secret-manager print-token"
```

On Unix, hcloud executes the string using `/bin/sh -c`; on Windows it uses
`cmd.exe /D /C`. Shell quoting, environment expansion, and pipelines follow that
shell's rules. Standard input and standard error are inherited directly from
hcloud, so a secret manager can prompt for a passphrase, PIN, or MFA code. Only
standard output is captured. No credential protocol or JSON response is needed.

Tokens are resolved lazily when authentication is needed and cached for the client
lifecycle. A nonempty `HCLOUD_TOKEN` overrides the command. Otherwise, a nonblank
`token_command` takes precedence over a stored context `token`. If `token_command`
is empty or contains only whitespace, the stored token is used. A command failure
or empty output stops authentication; it does not fall back to the stored token.

Creating, listing, renaming, or deleting contexts, inspecting configuration, and
displaying help do not execute token commands. Retrieved tokens are never written
back to configuration or included in configuration inspection, context-list
output, or debug traces. The command string itself is stored in the configuration
and can appear in context-list output; use it to reference secrets, rather than
embedding them directly. `--token-command` cannot be combined with
`--token-from-env` when creating a context.

### Passbolt example

With a Passbolt version that supports raw password output through
`get resource --field password --raw`, configure:

```toml
active_context = "production"

[[contexts]]
name = "production"
token_command = "passbolt get resource --id 01234567-89ab-cdef-0123-456789abcdef --field password --raw"
```

Passbolt manages its own server configuration and authentication. The resource's
password field contains the hcloud API token. Passbolt may prompt on stderr or
`/dev/tty`; its stdout must contain only the decrypted password. No adapter is
required. Other secret managers work the same way when they provide raw output.
