FROM golang:1.18-bullseye

WORKDIR /usr/src/app

COPY src/go.mod ./
COPY src/go.sum ./
COPY src/ .
RUN apt-get install -y tzdata

# DOWNLOAD DEPENDECIES
RUN go mod download
EXPOSE $PORT
RUN go install github.com/cosmtrek/air@latest
CMD [ "air" ]
