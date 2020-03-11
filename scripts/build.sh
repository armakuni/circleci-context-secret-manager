#!/usr/bin/env bash

set -euo pipefail

ARCH="amd64"
OS_TARGETS="darwin linux windows"

mkdir -p builds
for target in $OS_TARGETS; do
  GOOS="${target}" GOARCH="${ARCH}" go build -o "builds/ccsm-${target}-${ARCH}" ./cmd/ccsm
done
