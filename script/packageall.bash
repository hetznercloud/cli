#!/bin/bash -e

usage() {
  echo "usage: releaseall.bash RELEASE" >&2
  exit 2
}

release="$1"
[ -z "$release" ] && usage

cat script/variants.txt | while read os arch label; do
  echo $os-$arch
  script/build.bash $os $arch $release
  script/package.bash $os $arch $release
done
