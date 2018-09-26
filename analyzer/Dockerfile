FROM ubuntu:18.04
MAINTAINER ewaters@gmail.com

ADD FreeLing-4.1.tar.gz /
RUN cd /FreeLing-4.1 && \
    apt-get update -q && \
	apt-get install -y \
	    libboost-system1.65.1 libboost-program-options1.65.1 libboost-regex1.65.1 libboost-thread1.65.1 libboost-filesystem1.65.1 \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev libboost-filesystem-dev \
		libicu-dev zlib1g-dev git build-essential automake autoconf libtool cmake locales && \
	locale-gen en_US.UTF-8 && \
	mkdir build && \
	cd build && \
	cmake .. && \
	make install && \
	cd / && \
	rm -rf /FreeLing-4.1 && \
	apt-get purge -y \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev libboost-filesystem-dev \
		libicu-dev zlib1g-dev git build-essential automake autoconf libtool cmake && \
	apt-get autoremove -y && \
	apt-get clean -y && \
	rm -rf /var/lib/apt/lists/*
