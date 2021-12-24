# EXAMPLE USAGE:  docker run --rm codewars/python-runner run -l python -c "print 1+1 "
#                 docker run --rm codewars/python-runner run -l python3 -c "print(1+1)"

# Pull base image.
FROM ubuntu:bionic

ENV APP_DIR /app
ENV RUSTUP_HOME=/usr/local/rustup \
    CARGO_HOME=/usr/local/cargo \
    PATH=/usr/local/cargo/bin:$PATH

RUN apt-get update -qq && apt-get install -q -y --no-install-recommends \
        software-properties-common \
        build-essential

RUN apt-get install -y -q --no-install-recommends ca-certificates
RUN apt-get install -y -q --no-install-recommends curl
RUN rm -rf /var/lib/apt/lists/*

RUN mkdir -p $RUSTUP_HOME
RUN mkdir -p $CARGO_HOME

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y --no-modify-path --default-toolchain 1.15.1; \
    chmod -R a+w $RUSTUP_HOME $CARGO_HOME;

WORKDIR $APP_DIR
RUN mkdir -p $APP_DIR/src
COPY Cargo.toml .

RUN addgroup dockerusergroup
RUN chmod g+rwx $APP_DIR
RUN chgrp -R dockerusergroup $APP_DIR
RUN useradd -G dockerusergroup -d $APP_DIR -s /bin/bash dockeruser
RUN chown -R dockeruser:dockerusergroup $APP_DIR && chown dockeruser:dockerusergroup /usr/local/rustup && chown dockeruser:dockerusergroup /usr/local/rustup

USER dockeruser
ENV USER=dockeruser HOME=$APP_DIR