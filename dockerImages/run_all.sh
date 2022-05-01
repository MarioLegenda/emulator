#!/bin/bash

dir="/var/www/go/src/therebelsource/emulator/dockerImages"

cd "$dir/go_v1_17_6" && /usr/bin/docker image build -t go:go_v1_17_6 .
cd "$dir/node_latest" && /usr/bin/docker image build -t node:node_latest .
cd "$dir/node_latest" && /usr/bin/docker image build -t node:node_latest_esm .
cd "$dir/python2" && /usr/bin/docker image build -t python:python2 .
cd "$dir/python3" && /usr/bin/docker image build -t python:python3 .
cd "$dir/ruby" && /usr/bin/docker image build -t ruby:ruby .
cd "$dir/php7.4" && /usr/bin/docker image build -t php:php7.4 .
cd "$dir/rust" && /usr/bin/docker image build -t rust:rust .
cd "$dir/haskell" && /usr/bin/docker image build -t haskell:haskell .
cd "$dir/c" && /usr/bin/docker image build -t c:c .
cd "$dir/c++" && /usr/bin/docker image build -t c-plus:c-plus .
cd "$dir/c_sharp_mono" && /usr/bin/docker image build -t c_sharp_mono:c_sharp_mono .
