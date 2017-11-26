FROM ubuntu:latest
MAINTAINER ewaters@gmail.com

ADD FreeLing-master.tbz2 /
RUN cd /FreeLing-master && \
    apt-get update -q && \
	apt-get install -y \
		libboost-system1.58.0 libboost-program-options1.58.0 libboost-regex1.58.0 \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev \
		libicu-dev zlib1g-dev git build-essential automake autoconf libtool locales && \
	locale-gen en_US.UTF-8 && \
	autoreconf --install && \
	./configure && \
	make && make install && \
	cd / && \
	rm -rf /FreeLing-master && \
	apt-get purge -y \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev \
		libicu-dev zlib1g-dev git build-essential automake autoconf libtool && \
	apt-get autoremove -y && \
	apt-get clean -y && \
	rm -rf /var/lib/apt/lists/*
