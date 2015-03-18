#!/bin/bash

set -ex

DIR="$(cd "$(dirname "${0}")/.." && pwd)"
cd "${DIR}"

function mac_or_linux() {
  local uname_s="$(uname -s)"
  case "${uname_s}" in
    Darwin) echo "mac" ;;
    Linux) echo "linux" ;;
    *)
      echo "error: unknown result from uname -s: ${uname_s}" >&2
      return 1
  esac
}

HOST_TMPDIR=
case "$(mac_or_linux)" in
  mac) HOST_TMPDIR="${HOME}/tmp" ;;
  linux) HOST_TMPDIR="/tmp" ;;
esac
mkdir -p "${HOST_TMPDIR}"

GOOS=
case "${1}" in
  linux) GOOS=linux ;;
  darwin) GOOS=darwin ;;
  *)
    echo "error: must specify GOOS as first argument [linux,darwin]" >&2
    exit 1
esac

rm -rf "${HOST_TMPDIR}/goscm_build/${GOOS}"
docker run -v "${HOST_TMPDIR}:/tmp" pedge/goscmlib make "${GOOS}compile"
rm -rf "downloads/${GOOS}_amd64"
mkdir -p "downloads/${GOOS}_amd64"
cp "${HOST_TMPDIR}/goscm_build/${GOOS}/go-scm" "downloads/${GOOS}_amd64/go-scm"
echo "downloads/${GOOS}_amd64"
ls -lh "downloads/${GOOS}_amd64"
