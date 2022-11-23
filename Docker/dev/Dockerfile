FROM golang:1.18-bullseye

WORKDIR /usr/src/app

COPY src/go.mod ./
COPY src/go.sum ./
RUN apt-get install -y tzdata
RUN go mod download
EXPOSE $PORT
# COPY src/ .
RUN go install github.com/cosmtrek/air@latest
CMD [ "air" ]
