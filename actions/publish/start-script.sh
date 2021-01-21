#!/bin/bash
echo $AWS_ACCESS_KEY:$AWS_SECRET_ACCESS_KEY > /etc/passwd-s3fs && chmod 600 /etc/passwd-s3fs
s3fs $AWS_S3_BUCKET_NAME $AWS_S3_MOUNT_DIRECTORY
ls $AWS_S3_MOUNT_DIRECTORY