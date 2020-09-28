#!/bin/bash
set -e
#
#
# Create Win zip and push to GH Release asset
#
#
INTEGRATIONNAME=$1
ARCH=$2
TAG=$3
GITHUB_TOKEN=$4

echo "===> Creating zip nri-${INTEGRATIONNAME}-${ARCH}.${TAG:1}.zip"
mkdir -p zip/'New Relic'/'newrelic-infra'/'newrelic-integrations'/bin/
mkdir -p zip/'New Relic'/'newrelic-infra'/'integrations.d'/

cp dist/nri-redis-win_windows_$arch/nri-${INTEGRATIONNAME}.exe  zip/'New Relic'/'newrelic-infra'/'newrelic-integrations'/bin/
cp ${INTEGRATIONNAME}-definition.yml zip/'New Relic'/'newrelic-infra'/'newrelic-integrations'/
cp ${INTEGRATIONNAME}-config.yml.sample zip/'New Relic'/'newrelic-infra'/'integrations.d'/

cd zip
7z a -r nri-${INTEGRATIONNAME}-${ARCH}.${TAG:1}.zip .

echo "===> Pushing nri-${INTEGRATIONNAME}-${ARCH}.${TAG:1}.zip to GHA Release assets"
export $GITHUB_TOKEN
hub release edit -a nri-${INTEGRATIONNAME}-${ARCH}.${TAG:1}.zip -m "${TAG}" ${TAG}