#!/bin/bash
set -e
#
#
# Gets dist/zip_dirty created by Goreleaser and reorganize inside files
#
#

ARCH=$1
TAG=$2

hub release edit -a "build/package/windows/nri-${ARCH}-installer/bin/Release/nri-redis-${ARCH}.${TAG:1}.msi" -m ${TAG} ${TAG}