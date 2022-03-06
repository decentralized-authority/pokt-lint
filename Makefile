BUILD_DIR=./build
LAMBDA_PING_TEST_BINARY=LambdaPingTestHandler
LAMBDA_RELAY_TEST_BINARY=LambdaRelayTestHandler
PING_TEST_BINARY=pingtest
RELAY_TEST_BINARY=relaytest
LAMBDA_PING_TEST_TARGET=./cmd/LambdaPingTestHandler
LAMBDA_RELAY_TEST_TARGET=./cmd/LambdaRelayTestHandler
PING_TEST_TARGET=./cmd/pingtest
RELAY_TEST_TARGET=./cmd/relaytest

build-commands:
	go build -o ${BUILD_DIR}/${PING_TEST_BINARY} ${PING_TEST_TARGET}
	go build -o ${BUILD_DIR}/${RELAY_TEST_BINARY} ${RELAY_TEST_TARGET}

build-lambda:
	make lambda-pingtest
	make lambda-relaytest

lambda-pingtest:
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_PING_TEST_BINARY} ${LAMBDA_PING_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_PING_TEST_BINARY}

lambda-relaytest:
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_RELAY_TEST_BINARY} ${LAMBDA_RELAY_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_RELAY_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY}

test:
	go test -v ./...

clean:
	go clean
	rm -f ${LAMBDA_PING_TEST_BINARY} ${LAMBDA_RELAY_TEST_BINARY} ${PING_TEST_BINARY} ${RELAY_TEST_BINARY} ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY}.zip
