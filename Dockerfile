FROM golang:1.17-alpine

RUN apk add --no-cache make
RUN apk add --no-cache zip

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

VOLUME /app/build

ENTRYPOINT ["make"]

# build example:
# docker build -t pokt-lint .

# run example:
# docker run --rm -ti -v [path_to_host_build_output_dir]:/app/build pokt-lint build-lambda
