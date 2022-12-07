#!/bin/sh

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSWORD" --output ./.github/secrets/hcloud_cli.p12 ./.github/secrets/hcloud_cli.p12.gpg

security create-keychain -p "" build.keychain
# Use long timeout for keychain to avoid issues where codesign fails because the keychain is locked
# before it was used. Default timeout is 300s
security set-keychain-settings -u -t 3600 ~/Library/Keychains/build.keychain
security import ./.github/secrets/hcloud_cli.p12 -t agg -k ~/Library/Keychains/build.keychain -P "$CERT_PASSWORD" -A

security list-keychains -s ~/Library/Keychains/build.keychain
security default-keychain -s ~/Library/Keychains/build.keychain
security unlock-keychain -p "" ~/Library/Keychains/build.keychain

security set-key-partition-list -S apple-tool:,apple: -s -k "" ~/Library/Keychains/build.keychain
