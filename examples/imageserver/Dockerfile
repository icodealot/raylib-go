FROM golang:1.18-alpine
  
RUN apk update

RUN apk add \
    build-base \
    mesa-dev \
    mesa-osmesa

WORKDIR /api
COPY . .
RUN go get -u -v github.com/icodealot/raylib-go-headless@master
RUN go mod tidy

RUN go build -o /api/raylibserver

ENTRYPOINT /api/raylibserver