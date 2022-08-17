FROM golang:latest

WORKDIR /usr/src/app

ENV GO_ENV=production

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/

RUN go build -v -o /usr/local/bin/app ./src

CMD ["app"] 