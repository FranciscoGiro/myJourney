############################
# STEP 1 build executable binary
############################
FROM golang:1.17.0-stretch AS builder
# syntax=docker/dockerfile:1

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

COPY ./credentials.json /credentials.json
ENV GOOGLE_APPLICATION_CREDENTIALS /credentials.json

RUN go build -o /docker-gs-ping

EXPOSE 8080

CMD [ "/docker-gs-ping" ]