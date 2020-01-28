notarize {
  path = "/dist/hcloud-macos_darwin_amd64/hcloud"
  bundle_id = "com.hetzner-cloud.cli"
  staple = true
}

apple_id {
  username = "integrations@hetzner-cloud.de"
  password = "@env:HC_APPLE_DEVELOPER_PASSWORD"
}
