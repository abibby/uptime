# syntax=docker/dockerfile:1

FROM golang:1.17-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /uptime-cmd

ENV UPTIME_AGGREGATE_PATH=/uptime/aggregate
ENV UPTIME_LOG_PATH=/uptime/uptime

RUN mkdir /uptime

VOLUME /uptime

CMD [ "/uptime-cmd" ]
