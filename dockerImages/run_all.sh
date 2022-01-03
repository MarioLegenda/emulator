#!/bin/bash

dir="/home/mario/go-emulator/dockerImages"

cd "$dir/go_v1_14_2" && /usr/bin/docker image build -t go:go_v1_14_2 .
cd "$dir/node_latest" && /usr/bin/docker image build -t node:node_latest .
cd "$dir/node_v14_x" && /usr/bin/docker image build -t node:node_v14_x .
cd "$dir/python2" && /usr/bin/docker image build -t python:python2 .
cd "$dir/python3" && /usr/bin/docker image build -t python:python3 .
cd "$dir/ruby" && /usr/bin/docker image build -t ruby:ruby .
cd "$dir/php7.4" && /usr/bin/docker image build -t php:php7.4 .
cd "$dir/rust" && /usr/bin/docker image build -t rust:rust .
cd "$dir/haskell" && /usr/bin/docker image build -t haskell:haskell .
cd "$dir/c" && /usr/bin/docker image build -t c:c .
cd "$dir/c++" && /usr/bin/docker image build -t c-plus:c-plus .
