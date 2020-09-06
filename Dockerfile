FROM golang:1.14 AS builder

ENV APP_HOME /src
WORKDIR $APP_HOME

ENV GOPROXY https://goproxy.cn,direct
COPY go.* $APP_HOME/
RUN go mod download

COPY . .
RUN make build

# Runing environment
FROM debian:10

RUN echo \
    deb http://mirrors.aliyun.com/debian buster main \
    deb http://mirrors.aliyun.com/debian buster-updates main \
    deb http://mirrors.aliyun.com/debian-security buster/updates main \
    > /etc/apt/sources.list

RUN apt-get update \
    && apt-get install -y ca-certificates telnet procps curl

ENV APP_HOME /moreu
WORKDIR $APP_HOME

COPY --from=builder /src/build/bin/moreu /usr/local/bin