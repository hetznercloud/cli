#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

go install github.com/boumenot/gocover-cobertura
go install github.com/golang/mock/mockgen
