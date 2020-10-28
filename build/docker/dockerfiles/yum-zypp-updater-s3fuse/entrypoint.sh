#!/bin/bash
set -e
#
#
# Mount S3 with S3Fuse and update YUM or ZYPP repo.
#
#
OS_VERSIONS_LIST=($(echo $OS_VERSIONS | tr ',' "\n"))
ARCH_LIST=($(echo $ARCH | tr ',' "\n"))
SUFIX='1'

#################
#    S3 FUSE    #
#################

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

#######################################
#    UPLOAD TO S3, UPDATE METADATA    #
#######################################

echo "===> Importing GPG signature"
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 --decode | gpg --batch --import -

echo "===> Download RPM packages from GH"
for arch in "${ARCH_LIST[@]}"; do
  if [ $arch == 'x86_64' ]; then
    package_name="nri-${INTEGRATION}-${TAG:1}-${SUFIX}.${arch}.rpm"
  else
    package_name="nri-${INTEGRATION}-${TAG:1}-${arch}.rpm"
  fi
  echo "===> Download ${package_name} from GH"
  curl -SL https://github.com/${REPO_FULL_NAME}/releases/download/${TAG}/${package_name} -o ${package_name}
done

for arch in "${ARCH_LIST[@]}"; do
  for os_version in "${OS_VERSIONS_LIST[@]}"; do
    if [ $arch == 'x86_64' ]; then
      package_name="nri-${INTEGRATION}-${TAG:1}-${SUFIX}.${arch}.rpm"
    else
      package_name="nri-${INTEGRATION}-${TAG:1}-${arch}.rpm"
    fi
    LOCAL_REPO_PATH="${AWS_S3_MOUNTPOINT}${BASE_PATH}/${os_version}/${arch}"
    echo "===> Creating local directory if not exists ${LOCAL_REPO_PATH}/repodata"
    [ -d "${LOCAL_REPO_PATH}/repodata" ] || mkdir -p "${LOCAL_REPO_PATH}/repodata"
    echo "===> Uploading ${package_name} to S3 in ${BASE_PATH}/${os_version}/${arch}"
    cp ${package_name} ${LOCAL_REPO_PATH}
    echo "===> Delete and recreate metadata for ${package_name}"
    find ${LOCAL_REPO_PATH} -regex '^.*repodata' | xargs -n 1 rm -rf
    time createrepo --update -s sha "${LOCAL_REPO_PATH}"
    FILE="${LOCAL_REPO_PATH}/repodata/repomd.xml"
    while [ ! -f $FILE ];do
       echo "===> Waiting repomd.xml exists..."
    done
    echo "===> Updating GPG metadata dettached signature in ${BASE_PATH}/${os_version}/${arch}"
    gpg --batch --pinentry-mode=loopback --passphrase ${GPG_PASSPHRASE} --detach-sign --armor "${LOCAL_REPO_PATH}/repodata/repomd.xml"
  done
done

echo "===> umount s3 Fuse"
s3fs -o nonempty umount ${AWS_S3_MOUNTPOINT}