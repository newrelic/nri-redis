#!/bin/bash
set -e
#
#
# Gets dist/tarball_dirty created by Goreleaser and reorganize inside files
#
#

PROJECT_PATH=$1

for tarball_dirty in $(find dist -regex ".*_dirty\.\(tar.gz\)");do
  tarball=${tarball_dirty:5:${#tarball_dirty}-(5+13)} # Strips begining and end chars
  TARBALL_TMP="dist/tarball_temp"
  TARBALL_CONTENT_PATH="${TARBALL_TMP}/${tarball}_content"
  mkdir -p ${TARBALL_CONTENT_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/
  mkdir -p ${TARBALL_CONTENT_PATH}/etc/newrelic-infra/integrations.d/
  echo "===> Decompress ${tarball} in ${TARBALL_CONTENT_PATH}"
  tar -xvf ${tarball_dirty} -C ${TARBALL_CONTENT_PATH}

  echo "===> Move files inside ${tarball}"
  mv ${TARBALL_CONTENT_PATH}/nri-redis "${TARBALL_CONTENT_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/"
  mv ${TARBALL_CONTENT_PATH}/redis-definition.yml ${TARBALL_CONTENT_PATH}/var/db/newrelic-infra/newrelic-integrations/
  mv ${TARBALL_CONTENT_PATH}/redis-config.yml.sample ${TARBALL_CONTENT_PATH}/etc/newrelic-infra/integrations.d/

  echo "===> Creating tarball ${tarball}"
  cd ${TARBALL_CONTENT_PATH}
  tar -czvf ../${tarball} .
  cd $PROJECT_PATH
  echo "===> Moving tarball ${tarball}"
  mv "${TARBALL_TMP}/${tarball}" dist/
  echo "===> Cleaning dirty tarball ${tarball_dirty}"
  rm ${tarball_dirty}
done
