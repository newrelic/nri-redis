#!/bin/bash
set -e
#
#
# Download Windows zip to create msi
#
#

ARCH=$1
INTEGRATION=$2

zip_name="nri-${INTEGRATION}_windows_${TAG:1}_${ARCH}.zip"
wget -O ./dist --quiet https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${zip_name}

7z e "dist/${zip_name}" -o "dist/nri-${INTEGRATION}_windows_${TAG:1}_${ARCH}/"
cp "dist/nri-${INTEGRATION}_windows_${TAG:1}_${ARCH}/New Relic/newrelic-infra/newrelic-integrations/bin/nri-${INTEGRATION}.exe" dist/
