ARG VERSION=20.04
FROM ubuntu:$VERSION

ARG S3_MOUNT_DIRECTORY=/var/www/s3

# validate s3fs-fuse with the sec team
RUN apt-get update && \
  apt-get install -y s3fs
  #awscli

## S3fs-fuse credential config
WORKDIR $S3_MOUNT_DIRECTORY

## Entry Point
ADD start-script.sh /start-script.sh
RUN chmod 755 /start-script.sh
CMD ["/start-script.sh"]
