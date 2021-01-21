ARG VERSION=20.04
FROM ubuntu:$VERSION

ARG AWS_S3_MOUNT_DIRECTORY=/mnt/s3
ENV AWS_S3_MOUNT_DIRECTORY $AWS_S3_MOUNT_DIRECTORY

# validate s3fs-fuse with the sec team
RUN apt-get update && \
  apt-get install -y s3fs make curl

RUN mkdir $AWS_S3_MOUNT_DIRECTORY

WORKDIR /home/gha

ADD scripts ./scripts
ADD Makefile .

RUN mkdir ./assets

RUN chmod -R 755 ./scripts
CMD ["make"]
