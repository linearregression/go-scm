#!/bin/bash

set -ex

DIR="$(cd "$(dirname "${0}")/.." && pwd)"
cd "${DIR}"

GOOS=
case "${1}" in
  linux) GOOS=linux ;;
  darwin) GOOS=darwin ;;
  *)
    echo "error: must specify GOOS as first argument [linux,darwin]" >&2
    exit 1
esac

rm -rf "/tmp/goscm_build/${GOOS}"
mkdir -p "/tmp/goscm_build/${GOOS}"
CGO_ENABLED=0 GOOS=${GOOS} GOARCH=amd64 go build -a -compiler gc -installsuffix cgo -ldflags '-d -s -w' -o "/tmp/goscm_build/${GOOS}/scm-clone" cmd/scm-clone/main.go
CGO_ENABLED=0 GOOS=${GOOS} GOARCH=amd64 go build -a -compiler gc -installsuffix cgo -ldflags '-d -s -w' -o "/tmp/goscm_build/${GOOS}/scm-checkout" cmd/scm-checkout/main.go
mkdir -p "tmp/${GOOS}"
cp "/tmp/goscm_build/${GOOS}/scm-clone" "tmp/${GOOS}/scm-clone"
cp "/tmp/goscm_build/${GOOS}/scm-checkout" "tmp/${GOOS}/scm-checkout"
echo "tmp/${GOOS}"
ls -lh "tmp/${GOOS}"
