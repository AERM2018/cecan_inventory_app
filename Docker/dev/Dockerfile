FROM golang:1.18-alpine3.17

WORKDIR /usr/src/app

COPY src/go.mod ./
COPY src/go.sum ./
# COPY src/ .
RUN apt install -y tzdata

# DOWNLOAD DEPENDECIES
RUN go mod download
RUN go mod tidy
EXPOSE $PORT
RUN go install github.com/cosmtrek/air@latest
CMD [ "air" ]
