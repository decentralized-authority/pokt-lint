$(V).SILENT:

BUILD_DIR=./build
LAMBDA_PING_TEST_BINARY=LambdaPingTestHandler
LAMBDA_RELAY_TEST_BINARY=LambdaRelayTestHandler
PING_TEST_BINARY=pingtest
RELAY_TEST_BINARY=relaytest
LAMBDA_PING_TEST_TARGET=./cmd/LambdaPingTestHandler
LAMBDA_RELAY_TEST_TARGET=./cmd/LambdaRelayTestHandler
PING_TEST_TARGET=./cmd/pingtest
RELAY_TEST_TARGET=./cmd/relaytest

help: ## Show this help message.
	echo "usage: make [target] ..."
	echo
	echo "targets:"
	echo "-------"
	egrep '^(.+)\:\ ##\ (.+)' Makefile | column -t -c 2 -s ':#'

docserver: ## Run an interactive OpenAPI spec on port 3333
	docker-compose up -d 
	# Visit documentation at http://localhost:3333

docserver-stop: ## Stop the interactive spec
	docker-compose down

build-commands: ## compiles executables to ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${PING_TEST_BINARY} ${PING_TEST_TARGET}
	go build -o ${BUILD_DIR}/${RELAY_TEST_BINARY} ${RELAY_TEST_TARGET}

build-lambda: ## builds lambda function bundles in ${BUILD_DIR}
	make lambda-pingtest
	make lambda-relaytest

lambda-pingtest: ## builds the pingtest lambda function
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_PING_TEST_BINARY} ${LAMBDA_PING_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_PING_TEST_BINARY}

lambda-relaytest: ## builds the relaytest lambda function
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_RELAY_TEST_BINARY} ${LAMBDA_RELAY_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_RELAY_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY}

test: ## runs the unit tests
	go test -v ./...

clean: ## deletes build artifacts
	go clean
	rm -f ${LAMBDA_PING_TEST_BINARY} ${LAMBDA_RELAY_TEST_BINARY} ${PING_TEST_BINARY} ${RELAY_TEST_BINARY} ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY}.zip
