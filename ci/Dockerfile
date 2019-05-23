FROM golang:1.12
LABEL	maintainer="Mirko Boehm <mirko@endocode.com>"
ENV	LC_ALL C.UTF-8
ENV	LANG C.UTF-8
RUN	apt-get update
RUN	apt-get -yqq install python3 golang git libxml2-utils bash
ENV	SHELL /bin/bash
RUN	go get -u github.com/jstemmer/go-junit-report
ADD	. /shelldoc

