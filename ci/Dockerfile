FROM golang:1.12
LABEL	maintainer="Mirko Boehm <mboehm@openinventionnetwork.com>"
ENV	LC_ALL C.UTF-8
ENV	LANG C.UTF-8
RUN	apt-get update
RUN	apt-get -yqq install python3 golang git
RUN	go get github.com/endocode/shelldoc/cmd/shelldoc
RUN	apt-get -yqq install python3-pip
ADD	. /shelldoc

