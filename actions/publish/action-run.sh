#!/bin/bash

# build docker image form Dockerfile
docker build -t newrelic/infrastructure-publish-action -f ./actions/publish/Dockerfile ./actions/publish

# run docker container to perform all actions inside
docker run --rm --security-opt apparmor:unconfined --device /dev/fuse --cap-add SYS_ADMIN \
        -e AWS_SECRET_ACCESS_KEY -e AWS_ACCESS_KEY -e AWS_S3_BUCKET_NAME \
        newrelic/infrastructure-publish-action