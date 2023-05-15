# Changelog

## [1.34.0](https://github.com/hetznercloud/cli/compare/v1.33.2...v1.34.0) (2023-05-15)


### Features

* **servertype:** show included traffic ([#480](https://github.com/hetznercloud/cli/issues/480)) ([17c5f4f](https://github.com/hetznercloud/cli/commit/17c5f4f02f4753f6ce2b8e843725db9da1f78737))

## [1.33.2](https://github.com/hetznercloud/cli/compare/v1.33.1...v1.33.2) (2023-05-08)


### Bug Fixes

* **iso:** list only shows the first 50 results & missing field in json output ([#474](https://github.com/hetznercloud/cli/issues/474)) ([9d7c6a4](https://github.com/hetznercloud/cli/commit/9d7c6a416c33c98d30b6e5a0546a91ac25d5dced))

## v1.33.1

### What's Changed
* fix: crash on server create with missing server type by @apricote in https://github.com/hetznercloud/cli/pull/469


**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.33.0...v1.33.1

## v1.33.0

Affordable, sustainable & powerful! 🚀You can now get one of our Arm64 CAX servers to optimize your operations while minimizing your costs!
Discover Ampere’s efficient and robust Arm64 architecture and be ready to get blown away with its performance. 😎

Learn more: https://www.hetzner.com/news/arm64-cloud

### What's Changed
* test: fix gitlab test setup by @apricote in https://github.com/hetznercloud/cli/pull/466
* fix: send debug output to stderr by @apricote in https://github.com/hetznercloud/cli/pull/467
* feat: add support for ARM APIs by @apricote in https://github.com/hetznercloud/cli/pull/468


**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.32.0...v1.33.0

## v1.32.0

### Notable Changes

* Adding "loadbalancer" as an alias to the "load-balancer" command by @cedi in https://github.com/hetznercloud/cli/pull/424
* feat(primary-ip): add and remove labels by @apricote in https://github.com/hetznercloud/cli/pull/435
* feat(server): return password on rebuild by @apricote in https://github.com/hetznercloud/cli/pull/460
* fix(placement-group): invalid json response #464 by @apricote in https://github.com/hetznercloud/cli/pull/465

### All Changes

* Adding "loadbalancer" as an alias to the "load-balancer" command by @cedi in https://github.com/hetznercloud/cli/pull/424
* feat(primary-ip): add and remove labels by @apricote in https://github.com/hetznercloud/cli/pull/435
* chore: remove unused hcapi.CertificateClient by @samcday in https://github.com/hetznercloud/cli/pull/441
* chore: remove unused hcapi.PrimaryIPClient by @samcday in https://github.com/hetznercloud/cli/pull/442
* chore: remove unused hcapi.DataCenterClient by @samcday in https://github.com/hetznercloud/cli/pull/443
* chore: migrate hcapi.ISOClient usage to hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/444
* Adding a .devcontainer configuration for usage in GitHub Codespaces by @cedi in https://github.com/hetznercloud/cli/pull/419
* chore: replace hcapi.ImageClient usage with hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/445
* chore: replace hcapi.LocationClient usage with hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/446
* fix: improve unset version specifier by @apricote in https://github.com/hetznercloud/cli/pull/447
* Bump golang.org/x/net from 0.2.0 to 0.7.0 by @dependabot in https://github.com/hetznercloud/cli/pull/448
* chore: remove unused hcapi.PlacementGroupClient by @samcday in https://github.com/hetznercloud/cli/pull/450
* chore: migrate hcapi.SSHKeyClient usages to hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/449
* chore: migrate hcapi.VolumeClient usage to hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/451
* chore: replace hcapi.FloatingIPClient usages with hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/452
* chore: migrate hcapi.FirewallClient usages to hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/454
* chore: migrate hcapi.NetworkClient usages to hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/453
* chore: replace hcapi.LoadBalancerClient usages with hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/455
* chore: replace hcapi.ServerClient usages with hcapi2 by @samcday in https://github.com/hetznercloud/cli/pull/456
* chore(deps): update module github.com/hetznercloud/hcloud-go to v1.41.0 by @apricote in https://github.com/hetznercloud/cli/pull/459
* feat(server): return password on rebuild by @apricote in https://github.com/hetznercloud/cli/pull/460
* fix(placement-group): invalid json response #464 by @apricote in https://github.com/hetznercloud/cli/pull/465

### New Contributors
* @samcday made their first contribution in https://github.com/hetznercloud/cli/pull/441
* @dependabot made their first contribution in https://github.com/hetznercloud/cli/pull/448

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.31.1...v1.32.0

## v1.31.1

### What's Changed
* ci: fix issue where release pipeline fails and no assets are produced by @apricote in https://github.com/hetznercloud/cli/pull/430
* fix(ci): race-condition in signing macos binaries by @apricote in https://github.com/hetznercloud/cli/pull/433


**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.31.0...v1.31.1

## v1.31.0

### What's Changed
* server/list: Add missing PlacementGroup to JSON by @tomsiewert in https://github.com/hetznercloud/cli/pull/416
* Update the toml library to the latest version by @cedi in https://github.com/hetznercloud/cli/pull/422
* Adding an age column to the cli, closes #417 by @cedi in https://github.com/hetznercloud/cli/pull/420
* feat(completion): read network zones from API by @apricote in https://github.com/hetznercloud/cli/pull/426

### New Contributors
* @cedi made their first contribution in https://github.com/hetznercloud/cli/pull/422

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.30.4...v1.30.5

## v1.30.4

### What's Changed
* chore: update hcloud-go to v1.37.0 by @apricote in https://github.com/hetznercloud/cli/pull/413
* fix: primary-ip list returns max 50 items by @apricote in https://github.com/hetznercloud/cli/pull/415

### New Contributors
* @apricote made their first contribution in https://github.com/hetznercloud/cli/pull/414

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.30.3...v1.30.4

## v1.30.3

### What's Changed
* Fix hcloud server-type describe completion by @LKaemmerling in https://github.com/hetznercloud/cli/pull/407
* Improve hcloud server ssh command to use IPv6 automatically if no IPv… by @LKaemmerling in https://github.com/hetznercloud/cli/pull/406


**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.30.2...v1.30.3

## v1.30.2

### What's Changed
* Update Dependencies by @LKaemmerling in https://github.com/hetznercloud/cli/pull/402
* Fix primary-ip list  -o json by @LKaemmerling in https://github.com/hetznercloud/cli/pull/403


**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.30.1...v1.30.2

## v1.30.1

### What's Changed
* Fix hcloud server ssh with flexible network options by @LKaemmerling in https://github.com/hetznercloud/cli/pull/396


**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.30.0...v1.30.1

## v1.30.0

### What's Changed
* Add Alpine Linux to third-party packages by @firefly-cpp in https://github.com/hetznercloud/cli/pull/387
* Add Fedora to the list of third-party providers by @wULLSnpAXbWZGYDYyhWTKKspEQoaYxXyhoisqHf in https://github.com/hetznercloud/cli/pull/388
* fix(readme): correct messed-up columns ... by @wULLSnpAXbWZGYDYyhWTKKspEQoaYxXyhoisqHf in https://github.com/hetznercloud/cli/pull/389
* Remove freebsd64 rescue system type by @LKaemmerling in https://github.com/hetznercloud/cli/pull/391
* Remove Third-party packages Table by @LKaemmerling in https://github.com/hetznercloud/cli/pull/392
* Add Primary IP Support by @LKaemmerling in https://github.com/hetznercloud/cli/pull/393

### New Contributors
* @firefly-cpp made their first contribution in https://github.com/hetznercloud/cli/pull/387
* @wULLSnpAXbWZGYDYyhWTKKspEQoaYxXyhoisqHf made their first contribution in https://github.com/hetznercloud/cli/pull/388

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.29.5...v1.30.0

## v1.29.5

### What's Changed
* Fix: Use the correct object to return in case of created_from flag by @4ND3R50N in https://github.com/hetznercloud/cli/pull/385

### New Contributors
* @4ND3R50N made their first contribution in https://github.com/hetznercloud/cli/pull/385

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.29.4...v1.29.5

## v1.29.4

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.29.0...v1.29.4

## v1.29.1

### What's Changed
* Fix installation instructions by @fhofherr in https://github.com/hetznercloud/cli/pull/368
* Fix missing new line on hcloud describe command by @LKaemmerling in https://github.com/hetznercloud/cli/pull/380
* Use Go 1.18 for building & testing by @LKaemmerling in https://github.com/hetznercloud/cli/pull/381
* Trim and lowercase for column selectors  by @gadelkareem in https://github.com/hetznercloud/cli/pull/375

### New Contributors
* @gadelkareem made their first contribution in https://github.com/hetznercloud/cli/pull/375

**Full Changelog**: https://github.com/hetznercloud/cli/compare/v1.29.0...v1.29.1

## v1.29.0

- Add support for network zone `us-east`
- Build with Go 1.17

The binary for Apple Silicon is omitted for this release because of issues with the Apple Notary Service.

## v1.28.1

### Changelog

* 4410fb4 Fix panic on iso & location list as json (#361)
* 94b5d5f Move RDNS Commands to RDNS Client (#357)

## v1.28.0

### Changelog

* 3d7078a Add support for LB DNS PTRs (#355)
* eee45a9 Remove no longer used build scripts (#353)

## v1.27.0

### Changelog

* 4b8ed4d Placement groups (#352)

## v1.26.1

### Changelog

* 2ab6137 Fix firewall description
* 01180ad Update hcloud-go to 1.29.1

## v1.26.0

### Changelog

* b4c1d1b Add description field to firewall rules
* caa9bf2 Fix pagination of list commands (#347)

## v1.25.1

### Changelog

* 687f623 internal/cmd: Remove redundant DescribeJSON from JSONSchema (#345)


## v1.25.0

### Changelog

* 830d0bc Add support for App images (#344)
* ac23982 Fix changelog generation

## v1.25.0-alpha.1

### Changelog

This release contains a major refactoring of the code and is marked as alpha.  We recommend using the latest stable release, but feel free to test this release and report bugs if you find something.

## v1.24.0

### Changelog

* 08da869 Add Support for Firewall Protocol GRE & ESP (#331)

## v1.23.0

### Changelog

* 234dd6d Implement Firewall resource label selector (#328)
* 5ea977e5dda83022d701e056157f7e218c7674c6 Support getting Firewalls by label selector (#327)
* bb30002002cd2c8af6b20269eff549d09f7204a5 server: Add ability to get traffic as column (#325)
* 8d0f07e802cebf6df44daa3ad8933cebe489a8d1 firewall: Add empty slices for respective direction instead of nil-slices (#324)

## v1.22.1

### Changelog

* 4e97f5c Add Powershell completion (#316)
* b93bb4fe2716a34d79504e588d90f55dc8cf8ab9 Fix output option broke with last release (#315)

## v1.22.0

### Changelog

* 7969d5b Add support for managed certificates

## v1.21.1

### Changelog

* 5442833 The cli normalized the given CIDRs by default, so when a user entered 10.0.0.1/8 (as a sample) the cli normalized it to 10.0.0.0/8 silent. After this MR we now validate that the given IP is the start of the CIDR block (e.g. 10.0.0.0/8). (#304)

## v1.21.0

### Changelog

* 6c04c99 Specify timeout on release jobs
* afd597adb2e7bda63cd497546a7ecbb1186307cb  Implement Firewall support (#301)
* 67ba0adc61faf4ce4696626abb0c322029f6240d Update to Go 1.16 and support Apple Silicon (darwin/arm64) (#298)

## v1.20.0

### Changelog

* 178bf96 Add vswitch integration (#283)
* 9d209c0 Update to cobra 1.1.1 (#282)

## v1.19.1

### Changelog

## v1.19.0

### Changelog

* d5d2fec Update hcloud-go to 1.22.0 and expose correct disk size for resized without disk server (#269)
* 5049b00 Add handling for deprecated Images (#263)
* be48b5e Use go 1.15 (#267)
* ad3a564 Improve/Rewrite Shell completions (#266)

## v1.18.0

### Changelog

* 290c168 hcloud server describe use correct unit for traffic counter (#259)
* c1bd46c Implement Label Selector and IP target support (#258)
* d5a31ce Expose the new traffic fields and add load-balancer change-type command (#256)
*  Add pricing per location to load-balancer-type describe and server-type describe (#254)
* 5fc1464 Fix context list nil pointer when no active context was given. (#252)
* 8245b2f Add (required) to help text of all required args (#253)
* 03c3c82 Fix typos (#251)
* 506c1b1 added instructions for completion with antigen in zsh. (#240)
* 5d6f1bb Add command to request a VNC console (#238)

## v1.17.0

### Changelog

* 50a7de3 Add support for Load Balancers and Certificates (#245)
* 196557e Show Server Type CPU Type on server-type list, server-type describe and server describe (#244)
* b2d33f1 Allow the created field to be within the list responses (#237)

## v1.16.2

### Changelog

* 3bc0379 Fix completion of server name on hcloud server ssh (#233)
* cc8786c Update to go 1.14 (#234)
* 8c32195 Add missing labels to hcloud server list -o json response (#231)

## v1.16.1

### Changelog

* eef73ac Bugfix: Add private_net to server list json response (#229)

## v1.16.0

### Changelog

* 613eafc Add option to label servers, volumes, images, floating ips and ssh keys on creation (#227)
* 0ff7a1b Add JSON output option to all hcloud list commands (#225)

## v1.15.0

### Changelog

* 381f133 Switch Build and Release System to Github Actions (#223)
* 85e971e Add stale bot (#221)

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
