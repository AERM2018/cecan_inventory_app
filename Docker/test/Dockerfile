FROM golang:1.18-bullseye

WORKDIR /usr/src/app

COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download
EXPOSE 4000