source = ["./dist/hcloud-macos-build_darwin_amd64/hcloud"]
bundle_id = "cloud.hetzner.cli"

apple_id {
  username = "@env:HC_APPLE_DEVELOPER_USER"
  password = "@env:HC_APPLE_DEVELOPER_PASSWORD"
}

sign {
  application_identity = "@env:HC_APPLE_IDENTITY"
}

zip {
  output_path = "./dist/hcloud-macos-amd64.zip"
}
