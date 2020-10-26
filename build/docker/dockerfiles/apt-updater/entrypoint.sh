#!/bin/bash
set -e
#
#
# Uses github "depot" python script to update the APT repo in S3
#
#
SUFIX='1'
ARCH=( amd64 arm arm64 )
CODENAMES=( bionic buster jessie precise stretch trusty wheezy xenial )

echo "===> Installing Depot Pyhton script"
git clone ${DEPOT_REPO}
cd depot; python setup.py install

echo "===> Importing GPG signature and getting KeyId"
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 --decode | gpg --batch --import -
GPG_KEY_ID=$(gpg --list-secret-keys --keyid-format LONG | awk '/sec/{if (length($2) > 0) print $2}' | cut -d "/" -f2)
echo  "===> KEYiD: $GPG_KEY_ID"

mkdir -p /artifacts; cd /artifacts
DEB_PACKAGE="nri-${INTEGRATION}_${TAG:1}-${SUFIX}_${ARCH}.deb"
echo "===> Downloading ${DEB_PACKAGE} from GH"
curl -SL https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${DEB_PACKAGE} -o ${DEB_PACKAGE}

for arch in "${ARCH[@]}"; do
  for codename in "${CODENAMES[@]}"; do
     echo "===> Uploading to S3 $BASE_PATH ${DEB_PACKAGE} to component=main and codename=${codename}"
     POOL_PATH="pool/main/n/nri-${INTEGRATION}/${DEB_PACKAGE}"
     depot --storage="s3://${AWS_S3_BUCKET}/${BASE_PATH}" \
        --component=main \
        --codename=${codename} \
        --architecture=${arch} \
        --pool-path=${POOL_PATH} \
        --gpg-key ${GPG_KEY_ID} \
        --passphrase ${GPG_PASSPHRASE} \
        /artifacts/${DEB_PACKAGE} \
        --force
  done
done