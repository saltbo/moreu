FROM golang:1.14 AS builder

ENV APP_HOME /src
WORKDIR $APP_HOME

COPY go.* $APP_HOME/
ENV GOPROXY https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN make build

# Runing environment
FROM debian:10

RUN apt update && apt install -y procps

ENV APP_HOME /moreu
WORKDIR $APP_HOME

COPY deployments/ .
COPY --from=builder /src/build/bin/moreu /usr/local/bin