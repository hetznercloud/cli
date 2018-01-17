#!/usr/bin/env bash

set -e

# parse the current git commit hash
COMMIT=`git rev-parse HEAD`

# check if the current commit has a matching tag
TAG=$(git describe --exact-match --abbrev=0 --tags ${COMMIT} 2> /dev/null || true)
SUFFIX=''
DETAIL=''

# use the matching tag as the version, if available
if [ -z "$TAG" ]; then
  TAG=$(git describe --abbrev=0)
  COMMITS=$(git --no-pager log ${TAG}..HEAD --oneline)
  COMMIT_COUNT=$(echo -e "${COMMITS}" | wc -l)
  COMMIT_COUNT_PADDING=$(printf %03d $COMMIT_COUNT)
  SHORT_COMMIT_ID=$(git rev-parse --short HEAD)

  SUFFIX='-dev'
  DETAIL=".${COMMIT_COUNT_PADDING}.${SHORT_COMMIT_ID}"
fi

if [ -n "$(git diff --shortstat 2> /dev/null | tail -n1)" ]; then
  SUFFIX="-dirty"
fi


VERSION="${TAG}${SUFFIX}${DETAIL}"
echo $VERSION
