# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . /usr/src/app
RUN go mod download

COPY *.go /usr/src/app

RUN go build -o /usr/src/app/main
RUN chmod +x /usr/src/app/main

EXPOSE 8080 8080

CMD [ "/usr/src/app/main" ]