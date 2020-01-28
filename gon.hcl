source = ["./dist/hcloud-macos_darwin_amd64/hcloud"]
bundle_id = "com.hetzner-cloud.cli"
notarize {
  path = "./dist/hetzner-cloud-cli.dmg"
  staple = true
  bundle_id = "com.hetzner-cloud.cli"
}

apple_id {
  username = "integrations@hetzner-cloud.de"
  password = "@env:HC_APPLE_DEVELOPER_PASSWORD"
}

dmg {
  volume_name = "hetzner-cloud-cli"
  output_path = "./dist/hetzner-cloud-cli.dmg"
}
