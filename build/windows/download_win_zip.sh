#!/bin/bash
set -e
#
#
# Download Windows zip to create msi
#
#

ARCH=$1

zip_name="nri-redis_windows_${TAG:1}_${ARCH}.zip"
wget -O ./dist --quiet https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${zip_name}

7z e "dist/${zip_name}" -o ./dist
