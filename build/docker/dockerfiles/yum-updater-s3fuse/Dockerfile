FROM centos:8

ARG S3_FUSE_VERSION=v1.87

# Add S3Fuse https://github.com/s3fs-fuse/s3fs-fuse/wiki/Installation-Notes
RUN yum update -y
RUN yum install -y gcc libstdc++-devel gcc-c++ fuse fuse-devel curl-devel libxml2-devel mailcap git automake make openssl-devel jq wget rpm-sign \
    && git clone https://github.com/s3fs-fuse/s3fs-fuse \
    && cd s3fs-fuse/ \
    && git checkout tags/$S3_FUSE_VERSION \
    && ./autogen.sh \
    && export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig \
    && ./configure --prefix=/usr --with-openssl \
    && make \
    && make install

# Add Createrepo for repo metadata
RUN yum install -y curl createrepo

RUN yum clean all

ADD entrypoint.sh /
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]