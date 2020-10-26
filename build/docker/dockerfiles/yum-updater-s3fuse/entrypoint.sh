#!/bin/bash
set -e
#
#
# Mount S3 with S3Fuse and update YUM/ZYPP repo
#
#
OS_VERSIONS=( 5 6 7 8 )
ARCH=( x86_64 arm arm64)

[ "${DEBUG:-false}" == 'true' ] && { set -x; S3FS_DEBUG='-d -d'; }

# Defaults
: ${AWS_S3_AUTHFILE:='/root/.s3fs'}
: ${AWS_S3_MOUNTPOINT:='/mnt/repo'}
: ${AWS_S3_URL:='https://s3.amazonaws.com'}
: ${S3FS_ARGS:=''}

if [ ! -f "${AWS_S3_AUTHFILE}" ] && [ -z "$AWS_ACCESS_KEY_ID" ]; then
    echo "Error: AWS_ACCESS_KEY_ID not specified, or ${AWS_S3_AUTHFILE} not provided"
    exit 128
fi

if [ ! -f "${AWS_S3_AUTHFILE}" ] && [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
    echo "Error: AWS_SECRET_ACCESS_KEY not specified, or ${AWS_S3_AUTHFILE} not provided"
    exit 128
fi

# Write auth file if it does not exist
if [ ! -f "${AWS_S3_AUTHFILE}" ]; then
   echo "${AWS_ACCESS_KEY_ID}:${AWS_SECRET_ACCESS_KEY}" > ${AWS_S3_AUTHFILE}
   chmod 400 ${AWS_S3_AUTHFILE}
fi

echo "===> Mounting s3 in local docker with Fuse"
mkdir -p ${AWS_S3_MOUNTPOINT}
s3fs $S3FS_DEBUG $S3FS_ARGS -o passwd_file=${AWS_S3_AUTHFILE} -o url=${AWS_S3_URL} ${AWS_S3_BUCKET} ${AWS_S3_MOUNTPOINT}

echo "===> Importing GPG signature"
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 --decode | gpg --batch --import -

echo "===> Download packages from GH and uploading to S3"
for os_version in "${OS_VERSIONS[@]}"; do
  package_name="newrelic-infra-${TAG:1}.el${os_version}.${ARCH}.rpm"
  LOCAL_REPO_PATH="${AWS_S3_MOUNTPOINT}${BASE_PATH}/${os_version}/${ARCH}"

  echo "===> Downloading ${package_name} from GH"
  wget --quiet https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${package_name}

  echo "===>Creating local directory if not exists ${LOCAL_REPO_PATH}/repodata"
  [ -d "${LOCAL_REPO_PATH}/repodata" ] || mkdir -p "${LOCAL_REPO_PATH}/repodata"
  sleep 3

  echo "===> Uploading ${package_name} to S3 in ${BASE_PATH}/${os_version}/${ARCH}"
  cp ${package_name} ${LOCAL_REPO_PATH}

  echo "===> Updating metadata for $package_name"
  find ${LOCAL_REPO_PATH} -regex '^.*repodata' | xargs -n 1 rm -rf
  sleep 3
  time createrepo --update -s sha "${LOCAL_REPO_PATH}"
  FILE="${LOCAL_REPO_PATH}/repodata/repomd.xml"
  while [ ! -f $FILE ];do
     echo "===> Waiting repomd.xml exists..."
     sleep 5
  done

  echo "===>Updating GPG metadata dettached signature in ${BASE_PATH}/${os_version}/${ARCH}"
  gpg --batch --pinentry-mode=loopback --passphrase ${GPG_PASSPHRASE} --detach-sign --armor "${LOCAL_REPO_PATH}/repodata/repomd.xml"
done