FROM ubuntu:latest
MAINTAINER ewaters@gmail.com

RUN locale-gen en_US.UTF-8 && \
	apt-get update -q && \
	apt-get install -y \
		libboost-system1.58.0 libboost-program-options1.58.0 libboost-regex1.58.0 \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev \
		libicu-dev zlib1g-dev git build-essential automake autoconf libtool && \
	git clone https://github.com/TALP-UPC/FreeLing.git /fl && \
	cd /fl && \
		autoreconf --install && \
		./configure && \
		make && make install && \
	cd / && \
	rm -rf /fl && \
	apt-get purge -y \
		libboost-system-dev libboost-program-options-dev libboost-regex-dev libboost-thread-dev \
		libicu-dev zlib1g-dev git build-essential automake autoconf libtool && \
	apt-get autoremove -y
