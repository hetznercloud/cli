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
context-specific preferences. Contexts store a token or credential helper settings and can optionally have
additional preferences which take precedence over the globally set preferences.

However, a config file is not required. If no config file is found, the CLI will
use the default configuration. Overriding options using environment variables allows the
hcloud CLI to function in a stateless way. For example, setting `HCLOUD_TOKEN` is
already enough in many cases.

You can use the `hcloud config` command to manage your configuration, for example
to get, list, set and unset configuration options and preferences. You can view a list
of all available options and preferences [here](manual/hcloud_config.md#synopsis).

## Credential helpers

A context can retrieve its API token from an external credential helper instead of
storing the token in `cli.toml`. Helpers use the retrieval operation of the
[Git credential protocol](https://git-scm.com/docs/gitcredentials).

```sh
hcloud context create \
  --credential-helper /path/to/git-credential-helper \
  --credential-helper-arg=--some-option \
  --credential-url https://api.hetzner.cloud/production \
  production
```

`--credential-helper-arg` can be repeated; each occurrence supplies one argument,
including spaces or commas literally. hcloud runs the executable directly, with
these arguments followed by `get`. It does not invoke a shell, expand helper names
to `git-credential-*`, or read Git configuration. The executable can be an absolute
path or a name found on `PATH`.

The resulting context is stored as:

```toml
[[contexts]]
name = "production"
credential_helper = ["/path/to/git-credential-helper", "--some-option"]
credential_url = "https://api.hetzner.cloud/production"
```

`credential_url` identifies the secret to look up; it does not change the API
endpoint and hcloud does not connect to it. Its scheme, host (including an optional
port), path without the leading slash, and optional username are sent to the
helper. The path is always included when present, so different projects can select
different tokens. Renaming the context preserves its credential URL. URLs must not
contain passwords, query strings, fragments, or newline/NUL characters.

For the URL above, the helper receives the following input, terminated by a blank
line:

```text
protocol=https
host=api.hetzner.cloud
path=production

```

The helper must return a nonempty `password` field containing the API token:

```text
password=YOUR_API_TOKEN

```

hcloud retrieves credentials only when authentication is needed and caches the
result for the current process. Creating, listing, renaming, or deleting contexts,
inspecting configuration, and displaying help do not execute the helper. Retrieved
tokens are never written back to the configuration or included in configuration
inspection or context-list output, even with `--allow-sensitive`.

A nonempty `HCLOUD_TOKEN` overrides the helper without executing it. A context
cannot contain both a stored token and helper settings. Both helper settings are
required, and helper creation cannot be combined with `--token-from-env`.
An unsuccessful lookup stops the operation; there is no fallback to another
credential source. hcloud reports helper failures without displaying captured
stdout or stderr, since either may contain secrets. Manage and authenticate your
credential store separately; hcloud only invokes `get`, never `store` or `erase`.

### Passbolt example

The repository includes a Python 3
[Passbolt adapter](../../examples/credential-helper/passbolt.py) that translates
this protocol into `passbolt get resource --id <resource-id> --json`.
Install the adapter as an executable, install `passbolt` on `PATH`, and configure
Passbolt authentication separately. The example uses noninteractive Passbolt
authentication; it does not forward a terminal for passphrase or MFA prompts.
See the [Passbolt Go CLI documentation](https://www.passbolt.com/docs/user/quickstart/go-cli/).

Use the Passbolt resource ID as the credential URL path:

```sh
hcloud context create \
  --credential-helper /path/to/passbolt.py \
  --credential-url https://api.hetzner.cloud/01234567-89ab-cdef-0123-456789abcdef \
  production
```

The API token is stored in the resource's password field. The Passbolt server URL
and authentication belong to Passbolt's own configuration, not `credential_url`.
