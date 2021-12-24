FROM mhart/alpine-node:latest

ENV APP_DIR /app

WORKDIR $APP_DIR

ADD perms.sh $APP_DIR

RUN addgroup dockerusergroup
RUN chgrp -R dockerusergroup $APP_DIR
RUN adduser -D -G dockerusergroup -h $APP_DIR -s /bin/bash dockeruser
RUN chown -R dockeruser:dockerusergroup $APP_DIR

RUN chmod -R g+r $APP_DIR
RUN chown -R dockeruser:dockerusergroup /usr/bin/node && chown -R dockeruser /usr/bin/node

USER dockeruser
ENV USER=dockeruser HOME=$APP_DIR

RUN chmod +x $APP_DIR/perms.sh

CMD /app/perms.sh
RUN rm /app/perms.sh
