#!/bin/bash -e

usage() {
  echo "usage: build.bash OS ARCH RELEASE" >&2
  exit 2
}

os="$1"
[ -z "$os" ] && usage

arch="$2"
[ -z "$arch" ] && usage

release="$3"
[ -z "$release" ] && usage

GOOS=$os GOARCH=$arch go build -o ./dist/hcloud-$os-$arch-$release ./cmd/hcloud
