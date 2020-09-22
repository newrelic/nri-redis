#!/bin/bash
set -e
#
#
# Gets dist/tarball_fix created by Goreleaser and reorganize inside files
#
#
for tarball_fix in $(find dist -regex ".*\.\(tar.gz_fix\)");do
  tarball=${tarball_fix:5:${#tarball_fix}-9} # Strips begining and end chars
  TARBALL_PATH="dist/tarball_temp/${tarball}"
  mkdir -p ${TARBALL_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/
  mkdir -p ${TARBALL_PATH}/etc/newrelic-infra/integrations.d/
  echo "===> Decompress ${tarball} in ${TARBALL_PATH}"
  tar -xvf ${tarball_fix} -C ${TARBALL_PATH}

  echo "===> Move files inside ${tarball}"
  mv ${TARBALL_PATH}/nri-redis "${TARBALL_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/"
  mv ${TARBALL_PATH}/redis-definition.yml ${TARBALL_PATH}/var/db/newrelic-infra/newrelic-integrations/
  mv ${TARBALL_PATH}/redis-config.yml.sample ${TARBALL_PATH}/etc/newrelic-infra/integrations.d/

  echo "===> Creating tarball_fix ${tarball}"
  cd ${TARBALL_PATH}
  tar -czvf ${tarball} .
  mv ${tarball} ../../
done
