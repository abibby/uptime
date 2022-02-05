# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /uptime-cmd

ENV AGGREGATE_PATH=/uptime/aggregate
ENV LOG_PATH=/uptime/uptime

RUN mkdir /uptime

VOLUME /uptime

CMD [ "/uptime-cmd" ]
