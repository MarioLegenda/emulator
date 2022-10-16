FROM golang:1.18-alpine as golang_build

ENV APP_DIR /app

WORKDIR /app

RUN apk add --update docker openrc
RUN rc-update add docker boot
RUN apk add --no-cache --upgrade bash

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download && go mod tidy

COPY ./ .

RUN go build -o go_emulator_build

EXPOSE 3001

RUN ["chmod", "+x", "/app/entrypoint.sh"]

CMD ["/app/entrypoint.sh"]