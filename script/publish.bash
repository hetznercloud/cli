#!/bin/bash -e

usage() {
  echo "usage: publish.bash RELEASE" >&2
  exit 2
}

release="$1"
[ -z "$release" ] && usage

assets=()

while read -r os arch label; do
  if [ -f "dist/hcloud-$os-$arch-$release.tar.gz" ]; then
    asset="dist/hcloud-$os-$arch-$release.tar.gz"
  elif [ -f "dist/hcloud-$os-$arch-$release.zip" ]; then
    asset="dist/hcloud-$os-$arch-$release.zip"
   else
    echo "no asset found for $os-$arch" >&2
    exit 1
  fi
  assets+=(-a $asset)
done < script/variants.txt

hub release create -d -m "hcloud $release" ${assets[@]} $release
