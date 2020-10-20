#!/bin/bash
#
#
# Upload dist artifacts to GH Release assets
#
#
cd dist
for package in $(find  -regex ".*\.\(msi\|rpm\|deb\|zip\|tar.gz\)");do
  echo "===> Uploading to GH $TAG: ${package}"
  echo "===> Debugging"
  echo "GITHUB_TOKEN: $GITHUB_TOKEN"
  echo "GITHUB_USER: $GITHUB_USER"
  echo "===> END Debugging"
  hub release edit -a ${package} -m "${TAG}" ${TAG}
done
