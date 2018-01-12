#!/bin/bash -e

usage() {
  echo "usage: package.bash OS ARCH RELEASE" >&2
  exit 2
}

os="$1"
[ -z "$os" ] && usage

arch="$2"
[ -z "$arch" ] && usage

release="$3"
[ -z "$release" ] && usage

tmp="$(mktemp -d /tmp/hcloud-$os-$arch-$release.XXXXXXXXXX)"
trap "rm -rf $tmp" EXIT

mkdir $tmp/hcloud-$os-$arch-$release

if [ "$os" = "windows" ]; then
  cp dist/hcloud-$os-$arch-$release $tmp/hcloud-$os-$arch-$release/hcloud.exe
  cp LICENSE README.md $tmp/hcloud-$os-$arch-$release/
  (cd $tmp/ && zip - $(find hcloud-$os-$arch-$release -type f)) > dist/hcloud-$os-$arch-$release.zip
else
  mkdir $tmp/hcloud-$os-$arch-$release/bin
  cp dist/hcloud-$os-$arch-$release $tmp/hcloud-$os-$arch-$release/bin/hcloud
  cp LICENSE README.md $tmp/hcloud-$os-$arch-$release/
  (cd $tmp/ && tar czf - hcloud-$os-$arch-$release) > dist/hcloud-$os-$arch-$release.tar.gz
fi
