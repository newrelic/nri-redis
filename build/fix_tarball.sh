#!/bin/bash
set -e
#
#
# Creates Redis tarball and outputs it in repo root
#
#
LOCAL=true
VERSION=$1
ARCH=$2
SRC_TARBALL_PATH="dist/fix_tarball"
DST_TARBALL_PATH=$SRC_TARBALL_PATH

if $LOCAL;then
  SRC_TARBALL_PATH='dist'
fi

mkdir -p {"$DST_TARBALL_PATH"/var/db/newrelic-infra/newrelic-integrations/bin/,"$DST_TARBALL_PATH"/etc/newrelic-infra/integrations.d/}
tar -xvf $SRC_TARBALL_PATH/nri-redis_linux_$VERSION_$ARCH.tar.gz -C $DST_TARBALL_PATH
mv "$SRC_TARBALL_PATH"/nri-redis "$DST_TARBALL_PATH"/var/db/newrelic-infra/newrelic-integrations/bin/
mv "$SRC_TARBALL_PATH"/redis-definition.yml "$DST_TARBALL_PATH"/var/db/newrelic-infra/newrelic-integrations/
mv "$SRC_TARBALL_PATH"/redis-config.yml.sample "$DST_TARBALL_PATH"/etc/newrelic-infra/integrations.d/

tar -czvf nri-redis_linux_$VERSION_$ARCH.tar.gz -C $DST_TARBALL_PATH .
rm -rf $DST_TARBALL_PATH