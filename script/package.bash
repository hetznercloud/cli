#!/bin/bash -e

usage() {
  echo "usage: package.bash OS ARCH RELEASE" >&2
  exit 2
}

crlf() {
  sed $'s/$/\r/'
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

mkdir $tmp/hcloud-$os-$arch-$release/etc
go build -o $tmp/_hcloud ./cmd/hcloud
$tmp/_hcloud completion bash > $tmp/hcloud-$os-$arch-$release/etc/hcloud.bash_completion.sh
$tmp/_hcloud completion zsh > $tmp/hcloud-$os-$arch-$release/etc/hcloud.zsh_completion

if [ "$os" = "windows" ]; then
  cp dist/hcloud-$os-$arch-$release $tmp/hcloud-$os-$arch-$release/hcloud.exe
  cat LICENSE | crlf > $tmp/hcloud-$os-$arch-$release/LICENSE
  cat README.md | crlf > $tmp/hcloud-$os-$arch-$release/README.md
  (cd $tmp/ && zip - $(find hcloud-$os-$arch-$release -type f)) > dist/hcloud-$os-$arch-$release.zip
else
  mkdir $tmp/hcloud-$os-$arch-$release/bin
  cp dist/hcloud-$os-$arch-$release $tmp/hcloud-$os-$arch-$release/bin/hcloud
  cp LICENSE README.md $tmp/hcloud-$os-$arch-$release/
  (cd $tmp/ && tar czf - hcloud-$os-$arch-$release) > dist/hcloud-$os-$arch-$release.tar.gz
fi
