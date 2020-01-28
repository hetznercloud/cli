#!/bin/sh

gpg --quiet --batch --yes --decrypt --passphrase="$SECRETS_PASSWORD" --output ./.github/secrets/developerID_application.cer ./.github/secrets/developerID_application.cer.gpg

security create-keychain -p "" build.keychain
security import ./.github/secrets/developerID_application.cer -t agg -k ~/Library/Keychains/build.keychain -P "" -A

security list-keychains -s ~/Library/Keychains/build.keychain
security default-keychain -s ~/Library/Keychains/build.keychain
security unlock-keychain -p "" ~/Library/Keychains/build.keychain

security set-key-partition-list -S apple-tool:,apple: -s -k "" ~/Library/Keychains/build.keychain
