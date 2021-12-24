#!/bin/bash

cd go_v1_14_2 && /usr/bin/docker image build -t go:go_v1_14_2 .
cd ../node_latest && /usr/bin/docker image build -t node:node_latest .
cd ../node_v14_x && /usr/bin/docker image build -t node:node_v14_x .
cd ../python2 && /usr/bin/docker image build -t python:python2 .
cd ../python3 && /usr/bin/docker image build -t python:python3 .
cd ../ruby && /usr/bin/docker image build -t ruby:ruby .
cd ../php7.4 && /usr/bin/docker image build -t php:php7.4 .
cd ../rust && /usr/bin/docker image build -t rust:rust .
cd ../haskell && /usr/bin/docker image build -t haskell:haskell .
cd ../c && /usr/bin/docker image build -t c:c .
cd ../c++ && /usr/bin/docker image build -t c-plus:c-plus .
