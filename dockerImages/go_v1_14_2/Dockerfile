FROM ubuntu:bionic

ENV APP_DIR /app
ENV GOLANG_VERSION 1.14.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GO_BIN_PATH /usr/local/go/bin

RUN apt-get update -qq && apt-get install -q -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		pkg-config

RUN apt-get install -y -q --no-install-recommends curl
RUN apt-get install -y -q --no-install-recommends ca-certificates
RUN apt-get install -y -q --no-install-recommends tar
RUN apt-get install -y -q --no-install-recommends xz-utils
RUN rm -rf /var/lib/apt/lists/*

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /app
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go env -w GOPATH=$APP_DIR

WORKDIR /app

RUN addgroup dockerusergroup
RUN chmod g+rwx $APP_DIR
RUN chgrp -R dockerusergroup $APP_DIR
RUN useradd -G dockerusergroup -d $APP_DIR -s /bin/bash dockeruser
RUN chown -R dockeruser:dockerusergroup $APP_DIR && chown dockeruser:dockerusergroup /usr/local/go/bin && chown dockeruser:dockerusergroup /usr/local/go/bin

USER dockeruser

ENV USER=dockeruser HOME=$APP_DIR