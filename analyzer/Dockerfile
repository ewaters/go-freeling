FROM ubuntu:21.04
MAINTAINER ewaters@gmail.com

RUN \
    apt-get update -q && \
	apt-get install -y \
	    libboost-system1.74.0 libboost-program-options1.74.0 libboost-regex1.74.0 libboost-thread1.74.0 libboost-filesystem1.74.0 \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev libboost-filesystem-dev \
		libicu-dev zlib1g-dev build-essential automake autoconf libtool cmake locales curl && \
	locale-gen en_US.UTF-8 && \
	curl -SL https://github.com/TALP-UPC/FreeLing/releases/download/4.2/FreeLing-langs-src-4.2.tar.gz | tar -xzC / && \
	curl -SL https://github.com/TALP-UPC/FreeLing/releases/download/4.2/FreeLing-src-4.2.tar.gz | tar -xzC / && \
    cd /FreeLing-4.2 && \
	mkdir build && \
	cd build && \
	cmake .. && \
	make install && \
	cd / && \
	rm -rf /FreeLing-4.2 && \
	apt-get purge -y \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev libboost-filesystem-dev \
		libicu-dev zlib1g-dev build-essential automake autoconf libtool cmake curl && \
	apt-get autoremove -y && \
	apt-get clean -y && \
	rm -rf /var/lib/apt/lists/*
