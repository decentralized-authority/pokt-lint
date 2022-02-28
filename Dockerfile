FROM golang:1.17-alpine

RUN apk add --no-cache make
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN make build-commands