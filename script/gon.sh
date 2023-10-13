#!/usr/bin/env bash

set -eu

# Only sign on releasing
if [[ "${GITHUB_REF_TYPE:-}" != "tag" ]]; then
  exit 0
fi

BINARY_PATH="$1"

GON_CONFIG=$(mktemp gon_XXXX.json)
cleanup() {
  rm -f "$GON_CONFIG"
}
trap cleanup EXIT

printf '{
  "source": ["%s"],
  "bundle_id": "cloud.hetzner.cli",
  "apple_id": {
    "username": "integrations@hetzner-cloud.de",
    "password": "@env:HC_APPLE_DEVELOPER_PASSWORD"
  },
  "sign": {
    "application_identity": "Developer ID Application: Hetzner Cloud GmbH (4PM38G6W5R)"
  }
}' "$BINARY_PATH" > "$GON_CONFIG"

gon -log-level=debug "$GON_CONFIG"
