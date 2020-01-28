source = ["./dist/hcloud-macos_darwin_amd64/hcloud"]
bundle_id = "cloud.hetzner.cli"

apple_id {
  username = "integrations@hetzner-cloud.de"
  password = "@env:HC_APPLE_DEVELOPER_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Hetzner Cloud GmbH (4PM38G6W5R)"
}

dmg {
  volume_name = "hcloud-cli"
  output_path = "./dist/hcloud-cli.dmg"
}
