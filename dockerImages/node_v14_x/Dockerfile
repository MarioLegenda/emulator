FROM mhart/alpine-node:14.17.3

ENV APP_DIR /app

WORKDIR $APP_DIR

RUN adduser -D -h $APP_DIR -s /bin/bash dockeruser
RUN chown -R dockeruser:dockeruser $APP_DIR

RUN chmod -R g+r $APP_DIR
RUN chown -R dockeruser:dockeruser /usr/bin/node && chown -R dockeruser /usr/bin/node

USER dockeruser
ENV USER=dockeruser HOME=$APP_DIR