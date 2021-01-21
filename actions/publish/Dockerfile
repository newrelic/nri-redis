ARG VERSION=20.04
FROM ubuntu:$VERSION

ARG AWS_S3_MOUNT_DIRECTORY=/mnt/s3
ENV AWS_S3_MOUNT_DIRECTORY $AWS_S3_MOUNT_DIRECTORY

# validate s3fs-fuse with the sec team
RUN apt-get update && \
  apt-get install -y s3fs

WORKDIR $AWS_S3_MOUNT_DIRECTORY

## Entry Point
ADD start-script.sh /start-script.sh
RUN chmod 755 /start-script.sh
CMD ["/start-script.sh"]
