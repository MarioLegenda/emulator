FROM ruby:2.7-alpine3.11

ENV APP_DIR /app

WORKDIR $APP_DIR

RUN adduser -D -h $APP_DIR -s /bin/bash dockeruser
RUN chown -R dockeruser:dockeruser $APP_DIR

RUN chmod -R g+r $APP_DIR
RUN chown -R dockeruser:dockeruser /usr/local/bin/ruby && chown -R dockeruser /usr/local/bin/ruby

USER dockeruser
ENV USER=dockeruser HOME=$APP_DIR