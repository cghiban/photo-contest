FROM circleci/golang:1.16.5
MAINTAINER Daniel Jacobs <jacobs@cshl.edu>

ENV _APP Photo Contest With libvips
ENV _AUTHOR Daniel Jacobs
ENV _VERSION 0.0.1

WORKDIR ~/repo

RUN sudo mkdir -p /opt/bin
COPY install-vips.sh sleep.sh /opt/bin/

RUN sudo apt-get update \
	&& sudo apt-get install -y \
	&& sudo apt-get install -y libvips \
	&& sudo apt-get autoremove -y \
    && sudo apt-get clean \
    && sudo rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENTRYPOINT /opt/bin/sleep.sh
