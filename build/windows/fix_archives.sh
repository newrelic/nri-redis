#!/bin/bash
set -e
#
#
# Gets dist/zip_dirty created by Goreleaser and reorganize inside files
#
#

PROJECT_PATH=$1

for zip_dirty in $(find dist -regex ".*_dirty\.\(zip\)");do
  zip_file_name=${zip_dirty:5:${#zip_dirty}-(5+10)} # Strips begining and end chars
  ZIP_TMP="dist/zip_temp"
  ZIP_CONTENT_PATH="${ZIP_TMP}/${zip_file_name}_content"

  mkdir -p "${ZIP_CONTENT_PATH}"

  ls -la "${zip_dirty}"

  BIN_IN_ZIP_PATH="${ZIP_CONTENT_PATH}/New Relic/newrelic-infra/newrelic-integrations/bin/"
  CONF_IN_ZIP_PATH="${ZIP_CONTENT_PATH}/New Relic/newrelic-infra/integrations.d/"

  mkdir -p "${BIN_IN_ZIP_PATH}"
  mkdir -p "${CONF_IN_ZIP_PATH}"

  echo "===> Decompress ${zip_file_name} in ${ZIP_CONTENT_PATH}"
  unzip ${zip_dirty} -d ${ZIP_CONTENT_PATH}

  echo "===> Move files inside ${zip_file_name}"
  mv ${ZIP_CONTENT_PATH}/nri-redis.exe "${BIN_IN_ZIP_PATH}"
  mv ${ZIP_CONTENT_PATH}/redis-definition.yml "${CONF_IN_ZIP_PATH}"
  mv ${ZIP_CONTENT_PATH}/redis-config.yml.sample "${CONF_IN_ZIP_PATH}"

  echo "===> Creating zip ${zip_file_name}"
  cd "${ZIP_CONTENT_PATH}"
  zip -r ../${zip_file_name} .
  cd $PROJECT_PATH
  echo "===> Moving zip ${zip_file_name}"
  mv "${ZIP_TMP}/${zip_file_name}" dist/
  echo "===> Cleaning dirty zip ${zip_dirty}"
  rm "${zip_dirty}"
done