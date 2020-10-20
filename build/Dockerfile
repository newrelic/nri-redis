FROM golang:1.15.2-buster

RUN apt-get update \
    && apt-get -y install \
        rpm \
        gnupg2 \
        gpg-agent \
        debsigs \
        unzip \
        zip \
    && go get golang.org/dl/go1.9.7 \
    && /go/bin/go1.9.7 download

RUN curl -L https://github.com/cli/cli/releases/download/v1.1.0/gh_1.1.0_linux_amd64.deb -o gh_1.1.0_linux_amd64.deb
RUN dpkg -i gh_1.1.0_linux_amd64.deb
