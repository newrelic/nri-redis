#!/bin/bash

# build docker image form Dockerfile
docker build -t newrelic/infrastructure-publish-action -f ./actions/publish/Dockerfile ./actions/publish

# run docker container to perform all actions inside
docker run --rm --security-opt apparmor:unconfined \
        --device /dev/fuse \
        --cap-add SYS_ADMIN \
        -e AWS_SECRET_ACCESS_KEY \
        -e AWS_ACCESS_KEY \
        -e AWS_S3_BUCKET_NAME \
        -e REPO_NAME \
        -e APP_NAME \
        -e TAG \
        -e ARTIFACTS_DEST_FOLDER=$AWS_S3_MOUNT_DIRECTORY \
        -e ARTIFACTS_SRC_FOLDER=/home/gha/assets \
        -e UPLOADSCHEMA_FILE_PATH=/home/gha/schemas/ohi.yml \
        newrelic/infrastructure-publish-action
