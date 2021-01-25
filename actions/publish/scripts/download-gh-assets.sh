#!/bin/bash

download_pkg () {
  PKG_NAME=$1
  printf "downloading ${PKG_NAME} ... "

  set +e && curl -sS -L --fail -o ./assets/${PKG_NAME} "https://github.com/${GH_REPO}/releases/download/${GH_TAG}/${PKG_NAME}"
  test $? -eq 0 && echo "OK!"
}

download () {
  PKG_SCHEMA=$1
  read -r -a ARCHS <<< $2

  for arch in "${ARCHS[@]}"; do
    download_pkg "$(echo $PKG_SCHEMA | sed "s/ARCH/${arch}/g")"
  done
}
