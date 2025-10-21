package experimental

import "github.com/hetznercloud/cli/internal/cmd/base"

var DNS = base.ExperimentalWrapper("DNS API", "in beta", "https://docs.hetzner.cloud/changelog#2025-10-07-dns-beta")
var StorageBoxes = base.ExperimentalWrapper("Storage Box support", "experimental", "https://github.com/hetznercloud/cli/issues/1202")
