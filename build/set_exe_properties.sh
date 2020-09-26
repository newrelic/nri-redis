#!/bin/bash
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

AgentMajorVersion=$(echo $TAG | cut -d "." -f 1)
AgentMinorVersion=$(echo $TAG | cut -d "." -f 2)
AgentPatchVersion=$(echo $TAG | cut -d "." -f 3)
AgentBuildVersion='0'

sed \
  -e "s/{AgentMajorVersion}/$AgentMajorVersion/g" \
  -e "s/{AgentMinorVersion}/$AgentMinorVersion/g" \
  -e "s/{AgentPatchVersion}/$AgentPatchVersion/g" \
  -e "s/{AgentBuildVersion}/$AgentBuildVersion/g" versioninfo.json.template > src/versioninfo.json

echo "===> Checking versioninfo.json is in the path"

export PATH="$PATH:/go/bin"
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go generate github.com/newrelic/nri-redis/cmd/

