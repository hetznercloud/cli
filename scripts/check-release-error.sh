#!/usr/bin/env bash

ls -1 ./*.error || true

if [[ -f sign-and-notarize.error ]]; then exit 1; fi
