# Changes

## master

* Add `ssh` command for servers to start a new SSH connection to the server

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
