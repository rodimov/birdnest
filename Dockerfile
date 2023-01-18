FROM golang:1.18 AS builder

WORKDIR /usr/src/birdnest

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go mod tidy
RUN go build -v

FROM buildpack-deps:buster

WORKDIR /usr/local/bin

COPY --from=builder /usr/src/birdnest/birdnest .
COPY --from=builder /usr/src/birdnest/migrations ./migrations

ARG DBHOST
ARG DBPORT
ARG DBUSER
ARG DBPWD
ARG DBNAME
ARG DBDIALECT
ARG LOGSTORAGE
ARG SWAGHOST
ARG HOSTPORT

ENV DBHOST=$DBHOST
ENV DBPORT=$DBPORT
ENV DBUSER=$DBUSER
ENV DBPWD=$DBPWD
ENV DBNAME=$DBNAME
ENV DBDIALECT=$DBDIALECT
ENV LOGSTORAGE=$LOGSTORAGE
ENV SWAGHOST=$SWAGHOST
ENV HOSTPORT=$HOSTPORT

EXPOSE $HOSTPORT

CMD ["./birdnest"]
