FROM circleci/golang:1.16.5
MAINTAINER Daniel Jacobs <jacobs@cshl.edu>

ENV _APP Photo Contest With libvips
ENV _AUTHOR Daniel Jacobs
ENV _VERSION 0.0.1
ARG VIPS_VERSION="8.11.3"
ARG PKG_CONFIG_PATH="/home/circleci/vips/lib/pkgconfig"
ENV VIPS_VERSION $VIPS_VERSION
ENV PKG_CONFIG_PATH $PKG_CONFIG_PATH

WORKDIR ~/repo

RUN sudo mkdir -p /opt/bin
COPY install-vips.sh /opt/bin/

RUN sudo apt-get update && sudo apt-get install -y libglib2.0-dev libexpat1-dev libtiff5-dev libjpeg62-turbo-dev libgsf-1-dev libpng-dev gtk-doc-tools gobject-introspection && sudo /opt/bin/install-vips.sh
