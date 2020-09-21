#!/bin/bash
set -e
#
#
# Creates Redis tarball and outputs it in repo root
#
#
VERSION=$1
ARCH=$2
TARBALL_PATH='dist/fix_tarball'

mkdir -p ${TARBALL_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/
mkdir -p ${TARBALL_PATH}/etc/newrelic-infra/integrations.d/

echo "Decompression dist/nri-redis_linux_${VERSION}_${ARCH}.tar.gz in ${TARBALL_PATH}"
tar -xvf dist/nri-redis_linux_${VERSION}_${ARCH}.tar.gz -C ${TARBALL_PATH}


mv ${TARBALL_PATH}/nri-redis "${TARBALL_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/"
mv ${TARBALL_PATH}/redis-definition.yml ${TARBALL_PATH}/var/db/newrelic-infra/newrelic-integrations/
mv ${TARBALL_PATH}/redis-config.yml.sample ${TARBALL_PATH}/etc/newrelic-infra/integrations.d/

cd ${TARBALL_PATH}
tar -czvf nri-redis_linux_$VERSION_$ARCH.tar.gz . -C ../${TARBALL_PATH}

