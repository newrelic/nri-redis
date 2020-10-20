#!/bin/bash
set -e
#
#
# Upload dist artifacts to GH Release assets
#
#
cd dist
release_id=$(curl --header "authorization: Bearer $GITHUB_TOKEN" --url https://api.github.com/repos/${REPO_FULL_NAME}/releases/tags/${TAG} | jq --raw-output '.id' )

for filename in $(find  -regex ".*\.\(msi\|rpm\|deb\|zip\|tar.gz\)");do
  echo "===> Uploading to GH $TAG: ${filename}"
  curl -s \
       -H "Authorization: token $GITHUB_TOKEN" \
       -H "Content-Type: application/octet-stream" \
       --data-binary @$filename \
       "https://uploads.github.com/repos/${REPO_FULL_NAME}/releases/${release_id}/assets?name=$(basename $filename)"
done
