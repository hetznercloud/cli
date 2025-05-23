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
context-specific preferences. Contexts always store a token and can optionally have
additional preferences which take precedence over the globally set preferences.

However, a config file is not required. If no config file is found, the CLI will
use the default configuration. Overriding options using environment variables allows the
hcloud CLI to function in a stateless way. For example, setting `HCLOUD_TOKEN` is
already enough in many cases.

You can use the `hcloud config` command to manage your configuration, for example
to get, list, set and unset configuration options and preferences. You can view a list
of all available options and preferences [here](manual/hcloud_config.md#synopsis).
