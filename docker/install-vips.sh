#!/bin/bash 
# Inspired by https://github.com/libvips/php-vips/blob/master/install-vips.sh

version="8.11.3"
HOME="/home/circleci/"
vips_tarball=https://github.com/libvips/libvips/releases/download/v$version/vips-$version.tar.gz

set -e

rm -rf $HOME/vips
curl -Ls $vips_tarball | tar xz
cd vips-$version
./configure --prefix=$HOME/vips "$@"
make install
