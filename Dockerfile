FROM golang:1.17

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# build example:
# docker build -t pokt-lint .

# run example:
# docker run --rm -v $( pwd )/build:/app/build pokt-lint go build -o ./build/ ./cmd/...
