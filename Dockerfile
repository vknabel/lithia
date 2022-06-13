# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build ./app/lithia

##
## Deploy
##
FROM debian:stable-slim

WORKDIR /
COPY --from=build /app/lithia /bin/lithia
COPY ./stdlib /opt/lithia/stdlib
ENV LITHIA_STDLIB=/opt/lithia/stdlib

ENTRYPOINT ["/bin/lithia"]
