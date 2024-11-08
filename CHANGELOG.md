# Changelog

## [v1.49.0-rc.0](https://github.com/hetznercloud/cli/releases/tag/v1.49.0-rc.0)

### Features

- drop support for go 1.21 (#902)

### Bug Fixes

- macOS binary not notarized (#868)

## [1.48.0](https://github.com/hetznercloud/cli/compare/v1.47.0...v1.48.0) (2024-10-25)


### Features

* allow auto-completing context flag ([#861](https://github.com/hetznercloud/cli/issues/861)) ([7495b14](https://github.com/hetznercloud/cli/commit/7495b140e2359c47f5991bc2863674bb02dec9ea))
* **primary-ip:** allow setting auto-delete on create ([#860](https://github.com/hetznercloud/cli/issues/860)) ([b6aecb3](https://github.com/hetznercloud/cli/commit/b6aecb31893016e05070856c291704403b76a296)), closes [#855](https://github.com/hetznercloud/cli/issues/855)


### Bug Fixes

* **image:** only use numerical IDs for label auto-completion ([#865](https://github.com/hetznercloud/cli/issues/865)) ([1d10d92](https://github.com/hetznercloud/cli/commit/1d10d9284a7544971416c1efe8d2452a90b205e0))
* **primary-ip:** labels ignored in create command ([#885](https://github.com/hetznercloud/cli/issues/885)) ([7a17174](https://github.com/hetznercloud/cli/commit/7a1717412348f26bb14c27ee965f415044bb2b81)), closes [#883](https://github.com/hetznercloud/cli/issues/883)

## [1.47.0](https://github.com/hetznercloud/cli/compare/v1.46.0...v1.47.0) (2024-08-30)


### Features

* allow deleting rDNS entries ([#839](https://github.com/hetznercloud/cli/issues/839)) ([129f51d](https://github.com/hetznercloud/cli/commit/129f51de4cc76bb12d60a43ffd3068e1764ddbb9)), closes [#814](https://github.com/hetznercloud/cli/issues/814)
* **certificate:** allow retrying managed certificate issuance ([#847](https://github.com/hetznercloud/cli/issues/847)) ([0223f7d](https://github.com/hetznercloud/cli/commit/0223f7d7948d40b4dabe219c2578bc123a034bf2))
* improve datacenter describe server types readability ([#854](https://github.com/hetznercloud/cli/issues/854)) ([fc0002c](https://github.com/hetznercloud/cli/commit/fc0002c3fce9185a33bd485a8199200e4aec6776)), closes [#852](https://github.com/hetznercloud/cli/issues/852)


### Bug Fixes

* `quiet` option sometimes only working as a flag ([#846](https://github.com/hetznercloud/cli/issues/846)) ([f7c9ac6](https://github.com/hetznercloud/cli/commit/f7c9ac6241531c556544b5b398b474d05b67df9f))
* **image:** only accept numerical ids for update, delete, label ([#844](https://github.com/hetznercloud/cli/issues/844)) ([7355571](https://github.com/hetznercloud/cli/commit/735557183200eaeb4f5f8c586e5be84f245c52fc))
* **load-balancer:** show message if none of --server, --label-selector or --ip is set ([#849](https://github.com/hetznercloud/cli/issues/849)) ([2df45e0](https://github.com/hetznercloud/cli/commit/2df45e0082ec48faf40d64c6b313c37b4377d8c4))
* **primary-ip:** change protection commands do not allow protection levels ([#851](https://github.com/hetznercloud/cli/issues/851)) ([697a87a](https://github.com/hetznercloud/cli/commit/697a87ab2f0ddf6dbe53b3867a744230b2f8a952))

## [1.46.0](https://github.com/hetznercloud/cli/compare/v1.45.0...v1.46.0) (2024-07-25)

### API Changes for Traffic Prices and Server Type Included Traffic

There will be a breaking change in the API regarding Traffic Prices and Server Type Included Traffic on 2024-08-05. This release implements the necessary changes for the CLI.

Until now, the `hcloud server-type list` command had a column `traffic`. This is no longer possible, because there is no single included traffic for a server type, and instead this depends on the location. We have set the field to `-` and removed it from the default columns. The information is still available in `hcloud server-type describe <server-type-name>`, now in the "Pricing" section of the output.

This also affects the JSON and YAML output of both `hcloud server-type list` and `hcloud server-type describe` as we use the API schema for this.

You can learn more about this change in [our changelog](https://docs.hetzner.cloud/changelog#2024-07-25-cloud-api-returns-traffic-information-in-different-format).

### Features

* **load-balancer-type:** show traffic details for each location in describe output ([5c9380c](https://github.com/hetznercloud/cli/commit/5c9380c5b4d5792d259f13b6f44bf10f3d0c1b68))
* **server-type:** show traffic details for each location in describe output ([5c9380c](https://github.com/hetznercloud/cli/commit/5c9380c5b4d5792d259f13b6f44bf10f3d0c1b68))


### Bug Fixes

* ineffective poll interval option ([#822](https://github.com/hetznercloud/cli/issues/822)) ([41aaa56](https://github.com/hetznercloud/cli/commit/41aaa56bd4649d42054eabe48e14d9fdde42a8ae))
* printing for server create with multiple networks ([#824](https://github.com/hetznercloud/cli/issues/824)) ([f3fcbd2](https://github.com/hetznercloud/cli/commit/f3fcbd2378a18df09abc82c6b7ac9c8aa8a68302))
* **server-type:** remove deprecated traffic column ([5c9380c](https://github.com/hetznercloud/cli/commit/5c9380c5b4d5792d259f13b6f44bf10f3d0c1b68))

## [1.45.0](https://github.com/hetznercloud/cli/compare/v1.44.2...v1.45.0) (2024-07-17)


### Features

* allow configuring default sort columns for each supported resource ([#795](https://github.com/hetznercloud/cli/issues/795)) ([f6877a1](https://github.com/hetznercloud/cli/commit/f6877a166ea7602ac2f65b85d9dcae16690f49c6))
* better error messages on malformed option value ([#798](https://github.com/hetznercloud/cli/issues/798)) ([8c6fec9](https://github.com/hetznercloud/cli/commit/8c6fec9b4ef776c10d58efd86a5d2122b3a8be5c))


### Bug Fixes

* **config:** ordering of list option values not preserved ([#805](https://github.com/hetznercloud/cli/issues/805)) ([1ac27bf](https://github.com/hetznercloud/cli/commit/1ac27bfbcca1eb3dbb3778b1894badff3d7462f2))
* debug log is truncated if it already exists ([#793](https://github.com/hetznercloud/cli/issues/793)) ([c3d3a9f](https://github.com/hetznercloud/cli/commit/c3d3a9fa0e90340bc45c7b6ab1eeb0b318a12cf7))
* **firewall:** wrong wording when firewall is removed from resource ([#812](https://github.com/hetznercloud/cli/issues/812)) ([9017a65](https://github.com/hetznercloud/cli/commit/9017a652907b17cee0be9500bc36d455ee51c49f)), closes [#809](https://github.com/hetznercloud/cli/issues/809)

## [1.44.2](https://github.com/hetznercloud/cli/compare/v1.44.1...v1.44.2) (2024-07-01)


### Bug Fixes

* correctly print disk sizes on 32bit systems ([#800](https://github.com/hetznercloud/cli/issues/800)) ([3de8530](https://github.com/hetznercloud/cli/commit/3de8530dfc3b1e74f89bb69811ea5f872091ff18))

## [1.44.1](https://github.com/hetznercloud/cli/compare/v1.44.0...v1.44.1) (2024-06-26)


### Bug Fixes

* errors silently ignored when deleting resources ([#792](https://github.com/hetznercloud/cli/issues/792)) ([a7655f8](https://github.com/hetznercloud/cli/commit/a7655f8c42025b4576ee4556cb0ebbc1062671f5)), closes [#791](https://github.com/hetznercloud/cli/issues/791)

## [1.44.0](https://github.com/hetznercloud/cli/compare/v1.43.1...v1.44.0) (2024-06-20)


### Features

* delete multiple resources in parallel ([#761](https://github.com/hetznercloud/cli/issues/761)) ([f2fb321](https://github.com/hetznercloud/cli/commit/f2fb321a774419431bb851d8e92c3a4741bb9220))
* improve toml formatting ([#758](https://github.com/hetznercloud/cli/issues/758)) ([eacb7dd](https://github.com/hetznercloud/cli/commit/eacb7dd69c50a7d7469710b54adebc592702c2fd))
* **load-balancer:** allow specifying health check options in add-service ([#743](https://github.com/hetznercloud/cli/issues/743)) ([2cd08b2](https://github.com/hetznercloud/cli/commit/2cd08b2298882b5877c8ff522ddc4e2226ecbdb7)), closes [#742](https://github.com/hetznercloud/cli/issues/742)
* new action waiting progress  ([#749](https://github.com/hetznercloud/cli/issues/749)) ([9e30f3f](https://github.com/hetznercloud/cli/commit/9e30f3f69bc930e4b87be0cee8a096827b542d53))
* new configuration system, config subcommand ([#736](https://github.com/hetznercloud/cli/issues/736)) ([d1c6678](https://github.com/hetznercloud/cli/commit/d1c667838e045f9d93cc5d54fdd38bfb5f15ba17))
* **server-type:** add deprecated column to list command ([#780](https://github.com/hetznercloud/cli/issues/780)) ([906f864](https://github.com/hetznercloud/cli/commit/906f86436d1c6f7e9f0bebe43fcd26bd243f0b81))
* **server:** add default-ssh-keys option ([#759](https://github.com/hetznercloud/cli/issues/759)) ([9b34d26](https://github.com/hetznercloud/cli/commit/9b34d26295a675399741c7d85058233d3b9f0ee5))


### Bug Fixes

* **firewall:** 'create --rules-file' not working with outbound rules ([#752](https://github.com/hetznercloud/cli/issues/752)) ([2f2be32](https://github.com/hetznercloud/cli/commit/2f2be32a513f9e2e878553672189d24d098205bc)), closes [#750](https://github.com/hetznercloud/cli/issues/750)
* network list server count format ([#783](https://github.com/hetznercloud/cli/issues/783)) ([f69d261](https://github.com/hetznercloud/cli/commit/f69d26167100e83f1fa3ab0d658b624618e51c97))
* track progress if the terminal width allows it ([#768](https://github.com/hetznercloud/cli/issues/768)) ([069fffe](https://github.com/hetznercloud/cli/commit/069fffe48bb09f23edcade9e73f351a8b361a2a5)), closes [#767](https://github.com/hetznercloud/cli/issues/767)

## [1.43.1](https://github.com/hetznercloud/cli/compare/v1.43.0...v1.43.1) (2024-04-19)


### Bug Fixes

* allow passing multiple args to server ssh command ([#729](https://github.com/hetznercloud/cli/issues/729)) ([d98aee5](https://github.com/hetznercloud/cli/commit/d98aee5cf6fc94cbfd470c5515132adde0523e48)), closes [#728](https://github.com/hetznercloud/cli/issues/728)
* creating firewall with rules not working correctly ([64aeb23](https://github.com/hetznercloud/cli/commit/64aeb23f7f208d4fddd02d504fec9218c081dc04))
* **firewall:** missing error when not setting port for tcp rule ([#734](https://github.com/hetznercloud/cli/issues/734)) ([c5ad470](https://github.com/hetznercloud/cli/commit/c5ad47001ceb65157978eec2651ee6b198ac8233))
* **firewall:** setting rules not working correctly ([da1b684](https://github.com/hetznercloud/cli/commit/da1b6848ccf281cbf160cabf4cce2031245d52ab))

## [1.43.0](https://github.com/hetznercloud/cli/compare/v1.42.0...v1.43.0) (2024-04-03)


### Features

* allow deletion of multiple resources at once ([#719](https://github.com/hetznercloud/cli/issues/719)) ([3b896fe](https://github.com/hetznercloud/cli/commit/3b896fef82f50bc27351eec9e2da38032615c6e0))
* improve missing argument error messages ([#711](https://github.com/hetznercloud/cli/issues/711)) ([e7f9e74](https://github.com/hetznercloud/cli/commit/e7f9e74d22fa1c712ac8363d30efe0499916e859))
* **server:** allow JSON & YAML output in reset-password ([#716](https://github.com/hetznercloud/cli/issues/716)) ([373287b](https://github.com/hetznercloud/cli/commit/373287b452f2d3cf27c9d32648c3393ca45ded3a)), closes [#715](https://github.com/hetznercloud/cli/issues/715)


### Bug Fixes

* removing last rule from firewall fails with invalid_input error ([#696](https://github.com/hetznercloud/cli/issues/696)) ([acab17c](https://github.com/hetznercloud/cli/commit/acab17c361373e5bd3b689708195fa48f18dbafb))

## [1.42.0](https://github.com/hetznercloud/cli/compare/v1.41.1...v1.42.0) (2024-02-01)


### Features

* add global --quiet flag to hide non-error messages ([#656](https://github.com/hetznercloud/cli/issues/656)) ([25fcbbf](https://github.com/hetznercloud/cli/commit/25fcbbfe33680a61b9483956a30dd5403d471c3d)), closes [#644](https://github.com/hetznercloud/cli/issues/644)
* allow adding/removing multiple labels at once ([#665](https://github.com/hetznercloud/cli/issues/665)) ([919c446](https://github.com/hetznercloud/cli/commit/919c44649eec16e219954cbac4cc1b66a7f53380)), closes [#662](https://github.com/hetznercloud/cli/issues/662)
* group subcommands in command help ([#675](https://github.com/hetznercloud/cli/issues/675)) ([0cb271f](https://github.com/hetznercloud/cli/commit/0cb271f9859dab9bcf04e8b1be13ee4580356517))
* **server:** remove unsupported linux32 rescue type ([#679](https://github.com/hetznercloud/cli/issues/679)) ([5bb0350](https://github.com/hetznercloud/cli/commit/5bb03504474ac0ff35dd9fd550b6786a37eb39dd))


### Bug Fixes

* refetch after creating managed certificate ([#685](https://github.com/hetznercloud/cli/issues/685)) ([4864553](https://github.com/hetznercloud/cli/commit/4864553c53716a61a35b6c2b6ab0d1ed1156720a))
* **server:** fix typo in ip subcommand ([#678](https://github.com/hetznercloud/cli/issues/678)) ([c5e3f00](https://github.com/hetznercloud/cli/commit/c5e3f001236f2a2ef07b8181687e037d8f7b83a2))
* use --poll-interval flag ([#660](https://github.com/hetznercloud/cli/issues/660)) ([b9328a6](https://github.com/hetznercloud/cli/commit/b9328a6500837d92ac8632a6232f9d1f23338863))

## [1.41.1](https://github.com/hetznercloud/cli/compare/v1.41.0...v1.41.1) (2023-12-21)


### Bug Fixes

* allow multiple --type flags for metrics commands ([#649](https://github.com/hetznercloud/cli/issues/649)) ([69346f7](https://github.com/hetznercloud/cli/commit/69346f7a1b45480dd834f86788315f725516ae47)), closes [#642](https://github.com/hetznercloud/cli/issues/642)
* integer fields are displayed as floats when using create with o=yaml ([#651](https://github.com/hetznercloud/cli/issues/651)) ([2d42fcd](https://github.com/hetznercloud/cli/commit/2d42fcdc6095ed00cdda496ecc11f22803f73c46))
* integer fields are displayed as floats when using describe with -o=yaml ([#650](https://github.com/hetznercloud/cli/issues/650)) ([ee4c51e](https://github.com/hetznercloud/cli/commit/ee4c51ef2a3ced38988630645f0be9198f2bc222)), closes [#638](https://github.com/hetznercloud/cli/issues/638)
* wording when metrics are not yet available ([#639](https://github.com/hetznercloud/cli/issues/639)) ([9298c46](https://github.com/hetznercloud/cli/commit/9298c46b11bf3ab3b5a25eb571c95a89471b20f8))

## [1.41.0](https://github.com/hetznercloud/cli/compare/v1.40.0...v1.41.0) (2023-12-13)


### Features

* allow YAML output ([#628](https://github.com/hetznercloud/cli/issues/628)) ([87552d2](https://github.com/hetznercloud/cli/commit/87552d2a684f7147dc0528e32915fa96dd7332ad))
* **server:** allow to filter list by server status ([#629](https://github.com/hetznercloud/cli/issues/629)) ([1cb51f2](https://github.com/hetznercloud/cli/commit/1cb51f2f4148d9a0825c127910ddcb72ceac349f)), closes [#623](https://github.com/hetznercloud/cli/issues/623)


### Bug Fixes

* make JSON/YAML outputs consistent ([#622](https://github.com/hetznercloud/cli/issues/622)) ([43944b8](https://github.com/hetznercloud/cli/commit/43944b8561071883616f10cf03b6e714d3a4fa88))
* **server:** properly format request-console JSON output ([#626](https://github.com/hetznercloud/cli/issues/626)) ([40641b6](https://github.com/hetznercloud/cli/commit/40641b6f0ede6beb45ddbcd8362f6b0609e60eba))

## [1.40.0](https://github.com/hetznercloud/cli/compare/v1.39.0...v1.40.0) (2023-11-14)


### Features

* add warning messages to context commands ([#605](https://github.com/hetznercloud/cli/issues/605)) ([d84fe36](https://github.com/hetznercloud/cli/commit/d84fe3629f3879f157883bbcc21de618ae333e8b)), closes [#588](https://github.com/hetznercloud/cli/issues/588)
* allow JSON output on resource creation ([#609](https://github.com/hetznercloud/cli/issues/609)) ([d7241fe](https://github.com/hetznercloud/cli/commit/d7241fed634bdba2598197fe458a08205f6e0f36)), closes [#470](https://github.com/hetznercloud/cli/issues/470)
* **context:** ask to use HCLOUD_TOKEN when creating new context ([#582](https://github.com/hetznercloud/cli/issues/582)) ([54e582a](https://github.com/hetznercloud/cli/commit/54e582a7253274ef5b578b01fc9f992ef6eae226)), closes [#576](https://github.com/hetznercloud/cli/issues/576)
* pretty-print JSON output ([#610](https://github.com/hetznercloud/cli/issues/610)) ([ca4706f](https://github.com/hetznercloud/cli/commit/ca4706f118d4cfb559d84c7a726f99673d42f7d4))
* print out IPs after they were created ([#590](https://github.com/hetznercloud/cli/issues/590)) ([0e0eaa8](https://github.com/hetznercloud/cli/commit/0e0eaa8b58b184d26c8c1cb3b08b5b87bd5770f6))


### Bug Fixes

* "create" actions were not awaited in some cases ([#589](https://github.com/hetznercloud/cli/issues/589)) ([1c0e50a](https://github.com/hetznercloud/cli/commit/1c0e50ac88b78f0ae5b013e5a3a83fee9797d4df))
* **image:** backup is suggested as type in update command ([#603](https://github.com/hetznercloud/cli/issues/603)) ([6ea4695](https://github.com/hetznercloud/cli/commit/6ea469544a281d081fd5130ac5b7cd25df96ca18))
* inconsistent formatting in "describe" commands ([#593](https://github.com/hetznercloud/cli/issues/593)) ([b70dfe3](https://github.com/hetznercloud/cli/commit/b70dfe3613b75af0ef8b78cd01af8a2ae5dd47e5))
* **network:** create prints success to stderr instead of stdout ([#587](https://github.com/hetznercloud/cli/issues/587)) ([0a4c361](https://github.com/hetznercloud/cli/commit/0a4c3611b33dc8c1ae1621b03c45026fd0bd6950))
* **primary-ip:** list does now allow to output IP column ([#594](https://github.com/hetznercloud/cli/issues/594)) ([ee8e6ed](https://github.com/hetznercloud/cli/commit/ee8e6ed9a9e0ca2bc71dc04bbed8b80352a1a8dc)), closes [#591](https://github.com/hetznercloud/cli/issues/591)
* **primary-ip:** updating the name resets auto-delete ([#604](https://github.com/hetznercloud/cli/issues/604)) ([292f6c4](https://github.com/hetznercloud/cli/commit/292f6c40d33c8028184756318fe276b5b46c37c5))
* static resource list commands only return first 50 entries ([#592](https://github.com/hetznercloud/cli/issues/592)) ([7e511d9](https://github.com/hetznercloud/cli/commit/7e511d95669bade7057b0054c0bfafa57a6a1551))

## [1.39.0](https://github.com/hetznercloud/cli/compare/v1.38.3...v1.39.0) (2023-10-25)


### Features

* add --enable-protection flag to "create" commands ([#567](https://github.com/hetznercloud/cli/issues/567)) ([e313e69](https://github.com/hetznercloud/cli/commit/e313e6900f3fcf05eeace9af0c8697654b868df4))
* add "all list" command to list all resources in the project ([6d3b064](https://github.com/hetznercloud/cli/commit/6d3b064920f65807bccbf2f41f1acbc4836a760c))
* **iso:** allow to filter list by type (public, private) ([#573](https://github.com/hetznercloud/cli/issues/573)) ([140cbc3](https://github.com/hetznercloud/cli/commit/140cbc3931007e8b95e2e02d2bd9c20076da9d96))
* **primary-ip:** enable/disable-protection accept levels as arguments ([#564](https://github.com/hetznercloud/cli/issues/564)) ([b11e223](https://github.com/hetznercloud/cli/commit/b11e223c4ff51ebe46e452a10a22ca8ab002ac3b))
* **server:** add --enable-backup flag to "create" command ([#568](https://github.com/hetznercloud/cli/issues/568)) ([15adee0](https://github.com/hetznercloud/cli/commit/15adee05069e3470a9733c2cf95669436f88a253))
* **server:** add --wait flag to "shutdown" command ([#569](https://github.com/hetznercloud/cli/issues/569)) ([3ce048c](https://github.com/hetznercloud/cli/commit/3ce048cc576b21d7978daf308f48db75ebfc1f2f))


### Bug Fixes

* **floating-ip:** list command only returns first 50 entries ([#574](https://github.com/hetznercloud/cli/issues/574)) ([f3fa881](https://github.com/hetznercloud/cli/commit/f3fa8815dbec92d3f770dd2c441021aed5ce386b))
* **image:** list does not parse "type" flag correctly ([#578](https://github.com/hetznercloud/cli/issues/578)) ([9a0487a](https://github.com/hetznercloud/cli/commit/9a0487a5438e89feffe558f911522ec7b4daadf1))
* list outputs null instead of empty array when listing in JSON ([#579](https://github.com/hetznercloud/cli/issues/579)) ([93bed7e](https://github.com/hetznercloud/cli/commit/93bed7eb6b9c4d0f0b81f455c8f2ff2ba7e8e52b))

## [1.38.3](https://github.com/hetznercloud/cli/compare/v1.38.2...v1.38.3) (2023-10-16)


### Bug Fixes

* **build:** ensure signature is properly generated ([#562](https://github.com/hetznercloud/cli/issues/562)) ([77b313c](https://github.com/hetznercloud/cli/commit/77b313c4db3c4c707fd5ad454be79a3edf7e4d04))

## [1.38.2](https://github.com/hetznercloud/cli/compare/v1.38.2-rc.0...v1.38.2) (2023-10-13)


### Bug Fixes

* **build:** create release from previous candidates ([cf6eb47](https://github.com/hetznercloud/cli/commit/cf6eb472de8162c71f8de4355b714e6b0aa3a75f))

## [1.38.2-rc.0](https://github.com/hetznercloud/cli/compare/v1.38.1...v1.38.2-rc.0) (2023-10-13)


### Bug Fixes

* **build:** ensure unique tmp files for gon script ([#558](https://github.com/hetznercloud/cli/issues/558)) ([c20a78b](https://github.com/hetznercloud/cli/commit/c20a78b10c86747de5c50d117264666e6b5bb3c8))

## [1.38.1](https://github.com/hetznercloud/cli/compare/v1.38.0...v1.38.1) (2023-10-13)


### Bug Fixes

* **build:** goreleaser failed building binaries for release ([8e4cd29](https://github.com/hetznercloud/cli/commit/8e4cd2942e0b941ca0b9a61873214d9632614e76))

## [1.38.0](https://github.com/hetznercloud/cli/compare/v1.37.0...v1.38.0) (2023-10-12)


### Features

* build with Go 1.21 ([#543](https://github.com/hetznercloud/cli/issues/543)) ([368bfae](https://github.com/hetznercloud/cli/commit/368bfae953e074b4f6e81887bc437025b7dc0779))
* **iso:** support deprecation info API ([#555](https://github.com/hetznercloud/cli/issues/555)) ([2b0a0fa](https://github.com/hetznercloud/cli/commit/2b0a0fa47f01e5f22646e56840c7fa5663d2af6b))
* **load-balancer:** Add health status to list output ([#542](https://github.com/hetznercloud/cli/issues/542)) ([272cc63](https://github.com/hetznercloud/cli/commit/272cc635787a1ea09fb418fde5f5bba6252212d0))


### Bug Fixes

* typo in primary ipv6 not found error message ([#534](https://github.com/hetznercloud/cli/issues/534)) ([b9451f2](https://github.com/hetznercloud/cli/commit/b9451f2ac92bbcfb6201f6ca803c3fdc52cb557f))

## [1.37.0](https://github.com/hetznercloud/cli/compare/v1.36.0...v1.37.0) (2023-08-17)


### Features

* allow formatting a volume without automounting it ([#530](https://github.com/hetznercloud/cli/issues/530)) ([a435c9a](https://github.com/hetznercloud/cli/commit/a435c9a98a216eab3b1a2319092bd4a10a26cc9c))
* upgrade to hcloud-go v2 ([#512](https://github.com/hetznercloud/cli/issues/512)) ([e2df229](https://github.com/hetznercloud/cli/commit/e2df229c0f105c3138584424632a0a8ce3248e71))


### Bug Fixes

* make image subcommand descriptions consistent ([#519](https://github.com/hetznercloud/cli/issues/519)) ([34beff0](https://github.com/hetznercloud/cli/commit/34beff0910d63b9dae6a406c2076d3be4e23e760))
* **output:** ID column could not be selected ([#520](https://github.com/hetznercloud/cli/issues/520)) ([7d5594b](https://github.com/hetznercloud/cli/commit/7d5594bb29314b4eed5902302514fa73e1d9b445))
* **primary-ip:** assignee-id was not correctly passed when creating the IP ([#506](https://github.com/hetznercloud/cli/issues/506)) ([8c027b6](https://github.com/hetznercloud/cli/commit/8c027b65e6dd02b470f457c516aea3230e18b535))
* **server:** show actual progress on image-create ([a2f0874](https://github.com/hetznercloud/cli/commit/a2f0874af5e49d0c52df2dd5b2baebf39c7915e3))

## [1.36.0](https://github.com/hetznercloud/cli/compare/v1.35.0...v1.36.0) (2023-06-22)


### Features

* **network:** add support for exposing routes to vswitch connection ([#504](https://github.com/hetznercloud/cli/issues/504)) ([339cee9](https://github.com/hetznercloud/cli/commit/339cee9edb416b5055cf2d401124d2b9efe4ab1d))

## [1.35.0](https://github.com/hetznercloud/cli/compare/v1.34.1...v1.35.0) (2023-06-13)


### Features

* show server-type deprecation warnings ([#490](https://github.com/hetznercloud/cli/issues/490)) ([c5c0527](https://github.com/hetznercloud/cli/commit/c5c052732f0e87f7040640e20f372d8b2c2ba315))

## [1.34.1](https://github.com/hetznercloud/cli/compare/v1.34.0...v1.34.1) (2023-06-01)


### Bug Fixes

* **server:** wait for delete before returning ([#482](https://github.com/hetznercloud/cli/issues/482)) ([62cb07f](https://github.com/hetznercloud/cli/commit/62cb07f5aa6938cbdb066113e42672f16e882287))

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
