FROM ubuntu:bionic

ENV APP_DIR /app
ENV GHC_PATH /opt/ghc/bin
ENV CABAL_PATH /opt/cabal/bin

RUN apt-get update -qq && apt-get install -q -y --no-install-recommends \
        software-properties-common \
		g++ \
		gcc \
		libc6-dev \
		pkg-config

RUN apt-get install -y -q --no-install-recommends curl
RUN apt-get install -y -q --no-install-recommends ca-certificates
RUN apt-get install -y -q --no-install-recommends tar
RUN apt-get install -y -q --no-install-recommends xz-utils
RUN add-apt-repository -y ppa:hvr/ghc
RUN apt-get update
RUN apt-get install -y cabal-install-3.2 ghc-8.8.3

ENV PATH /usr/local/bin:$GHC_PATH:$CABAL_PATH:$PATH

WORKDIR /app

RUN addgroup dockerusergroup
RUN chgrp -R dockerusergroup $APP_DIR
RUN useradd -G dockerusergroup -d $APP_DIR -s /bin/bash dockeruser
RUN chown -R dockeruser:dockerusergroup $APP_DIR && chown dockeruser:dockerusergroup $GHC_PATH && chown dockeruser:dockerusergroup $GHC_PATH
RUN chown -R dockeruser:dockerusergroup $APP_DIR && chown dockeruser:dockerusergroup $CABAL_PATH && chown dockeruser:dockerusergroup $CABAL_PATH
RUN chmod -R 777 $APP_DIR

USER dockeruser:dockerusergroup

ENV USER=dockeruser HOME=$APP_DIR