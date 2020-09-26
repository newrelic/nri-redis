#!/bin/bash
set -e
#
#
# Create the metadata for the exe's files, called by .goreleser as a hook in the build section
#
#
TAG=$1

if [ -n "$1" ]; then
  echo "===> Tag is ${TAG}"
else
  echo "===> Tag not specified will be 0.0.0"
  TAG='0.0.0'
fi

MajorVersion=$(echo ${TAG:1} | cut -d "." -f 1)
MinorVersion=$(echo ${TAG:1} | cut -d "." -f 2)
PatchVersion=$(echo ${TAG:1} | cut -d "." -f 3)
BuildVersion='0'

sed \
  -e "s/{MajorVersion}/$MajorVersion/g" \
  -e "s/{MinorVersion}/$MinorVersion/g" \
  -e "s/{PatchVersion}/$PatchVersion/g" \
  -e "s/{BuildVersion}/$BuildVersion/g" versioninfo.json.template > src/versioninfo.json

echo "===> Checking versioninfo.json is in the path"
cat src/versioninfo.json

export PATH="$PATH:/go/bin"
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go generate github.com/newrelic/nri-redis/cmd/

