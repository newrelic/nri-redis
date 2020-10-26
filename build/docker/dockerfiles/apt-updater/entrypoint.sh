#!/bin/bash
set -e
#
#
# Uses github "depot" python script to update the APT repo in S3
#
#
if [ $PIPELINE_ACTION == 'prereleased' ]; then
  CODENAMES=( bionic )
  BOOT=( systemd )
fi
if [ $PIPELINE_ACTION == 'released' ]; then
  CODENAMES=( bionic buster jessie precise stretch trusty wheezy xenial )
  BOOT=( systemd upstart sysv )
fi

echo "===> Installing Depot Pyhton script"
git clone $DEPOT_REPO
cd depot; python setup.py install

echo "===> Importing GPG signature and getting KeyId"
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 --decode | gpg --batch --import -
GPG_KEY_ID=$(gpg --list-secret-keys --keyid-format LONG | awk '/sec/{if (length($2) > 0) print $2}' | cut -d "/" -f2)

echo "===> Downloading DEB packages from GH"
mkdir -p /artifacts; cd /artifacts
for boot in "${BOOT[@]}"; do
  echo "===> Downloading newrelic_infra_${boot}_${TAG:1}_amd64.deb from GH"
  DEB_PACKAGE="newrelic-infra_${boot}_${TAG:1}_amd64.deb"
  curl -SL https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${DEB_PACKAGE} -o ${DEB_PACKAGE}
done

for codename in "${CODENAMES[@]}"; do
  for boot in "${BOOT[@]}"; do
   echo "==> Release: Uploading to S3 newrelic-infra_${boot}_${TAG:1}_amd64.deb to component=main and codename=${codename}"
   DEB_PACKAGE="newrelic-infra_${boot}_${TAG:1}_amd64.deb"
   POOL_PATH="pool/main/n/newrelic-infra/${DEB_PACKAGE}"
   depot --storage=${AWS_S3_REPO_URL}/${BASE_PATH} \
      --component=main \
      --codename=${codename} \
      --pool-path=${POOL_PATH} \
      --gpg-key ${GPG_KEY_ID} \
      --passphrase ${GPG_APT_PASSPHRASE} \
      /artifacts/${DEB_PACKAGE} \
      --force
  done
done