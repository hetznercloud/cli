source = ["./dist/hcloud-macos-build_darwin_arm64/hcloud"]
bundle_id = "cloud.hetzner.cli"

apple_id {
  username = "integrations@hetzner-cloud.de"
  password = "@env:HC_APPLE_DEVELOPER_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Hetzner Cloud GmbH (4PM38G6W5R)"
}

zip {
  output_path = "./dist/hcloud-macos-arm64.zip"
}
