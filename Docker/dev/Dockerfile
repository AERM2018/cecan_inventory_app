FROM golang:1.16.15

WORKDIR /usr/src/app

COPY src/go.mod ./
COPY src/go.sum ./
COPY src/ .
RUN apt install -y tzdata

# DOWNLOAD DEPENDECIES
RUN go mod download
EXPOSE $PORT
RUN go install github.com/cosmtrek/air@latest
CMD [ "air" ]
