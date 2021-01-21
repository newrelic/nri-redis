#!/bin/bash

GH_REPO=newrelic/nri-redis
GH_TAG=v1.6.0
PKG_NAME=nri-redis
PKG_VERSION=1.6.0

#RPM_ARCHS=(386 1.x86_64 arm arm64)
#RPM='${PKG_NAME}-${PKG_VERSION}-${ARCH}.rpm'
#
#DEB_ARCHS=(386 amd64 arm arm64)
#DEB='${PKG_NAME}_${PKG_VERSION}-1_${ARCH}.deb'

#TAR_ARCHS=(386 amd64 arm arm64)
#TAR='${PKG_NAME}_linux_${PKG_VERSION}_${ARCH}.tar.gz'

download_pkg () {
  PKG_NAME=$1
  set +e && curl -sS -L -o ${PKG_NAME} "https://github.com/${GH_REPO}/releases/download/${GH_TAG}/${PKG_NAME}"
}

download () {
  ARCHS=$1
  PKG_SCHEMA=$2

  for arch in "${ARCHS[@]}"; do
    download_pkg "${PKG_NAME}-${arch}.${PKG_VERSION}.msi"
  done
}

WIN_ARCHS=(386 amd64)
MSI='${PKG_NAME}-${ARCH}.${PKG_VERSION}.msi'
download $WIN_ARCHS $MSI

ZIP='${PKG_NAME}-${ARCH}.${PKG_VERSION}.zip'
download $WIN_ARCHS $ZIP
