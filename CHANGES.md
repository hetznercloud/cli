# Changes

As of v1.15.0 we expose the Changelog only in the [Github Releases section](https://github.com/hetznercloud/cli/releases). 

## v1.14.0

* Expose server’s MAC address in networks on `hcloud server describe`
* Add support for names to Floating IP commands
* Make `--ip-range` on `hcloud network add-subnet` optional
* Add debug mode (use `HCLOUD_DEBUG` and `HCLOUD_DEBUG_FILE`)
* Add `hcloud server ip` command
* Expose `Created` on `hcloud floating-ip|image|ssh-key|volume describe`
* Refactor progressbar and add loading animation for running actions

## v1.13.0

* Show server name instead of ID on `hcloud floating-ip|volume|image list`
* Add support for networks

## v1.12.0

* Add support for executing commands via `hcloud server ssh <server> <command>`
* Make overriding context via `HCLOUD_CONTEXT` work
* Add support for JSON and Go template output
* Add support for multiple user data files
* Add length validation for API token on `hcloud context create`
* Add `active` column to context list on `hcloud context list`

## v1.11.0

* Add support for automounting and formatting volumes

## v1.10.0

* Fix creating a volume when server is specified by its name
* Deprecate and ignore the `--window` flag on `hcloud server enable-backup`
* Add output columns `type|labels|volumes|protection` to `hcloud server list`
* Add output columns `labels|protection` to `hcloud volume list`
* Add output column `labels` to `hcloud image list`
* Add output column `labels` to `hcloud floating-ip list`
* Add output column `labels` to `hcloud ssh-key list`

## v1.9.1

* Fix formatting issue on `hcloud volume list` and `hcloud volume describe`

## v1.9.0

* Add support for volumes
* Add `--start-after-create` flag to `hcloud server create` command

## v1.8.0

* Add `hcloud ssh-key update` command
* Add `-u/--user` and `-p/--port` flags to `hcloud server ssh` command
* Add `hcloud server set-rdns` command
* Add `hcloud floating-ip set-rdns` command

## v1.7.0

* Add type filter flag `-t` / `--type` to `image list` command
* Expose labels of servers, Floating IPs, images, and SSH Keys
* Add `hcloud {server|ssh-key|image|floating-ip} {add-label|remove-label}` commands

## v1.6.1

* Fix invalid formatting of integers in `hcloud * list` commands

## v1.6.0

* Show IP address upon creating a server
* Add `--poll-interval` flag for controlling the polling interval (for example for action progress updates)

## v1.5.0

* Add `hcloud server ssh` command to open an SSH connection to the server

## v1.4.0

* Document `-o` flag for controlling output formatting
* Add commands `enable-protection` and `disable-protection` for
  images, Floating IPs, and servers

## v1.3.2

* Show progress for every action
* Show datacenter in `server list` and `server describe`

## v1.3.1

* Only poll action progress every 500ms (instead of every 100ms)
* Document `HCLOUD_TOKEN` and make it work when there is no active context

## v1.3.0

* Print dates in local time
* Do not echo token when creating a context
* Add `--user-data-from-file` flag to `hcloud server create` command

## v1.2.0

* Update hcloud library to v1.2.0 fixing rate limit check

## v1.1.0

* Show image information in `hcloud server describe`
* Auto-activate created context on `hcloud context create`
* Fix `hcloud version` not showing correct version
