FROM ubuntu:21.04
MAINTAINER ewaters@gmail.com

RUN apt-get update -q && \
    apt-get install -y \
        git golang && \
    cd / && \
    git clone https://github.com/ewaters/go-freeling.git && \
    cd /go-freeling/proxy && \
    go build && \
    mv proxy /usr/local/bin/freeling-proxy && \
    cd / && \
    rm -rf /go-freeling && \
    apt-get purge -y \
        git golang && \
    apt-get autoremove -y && \
    apt-get clean -y && \
    rm -rf /var/lib/apt/lists/*
