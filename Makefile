$(V).SILENT:

BUILD_DIR=${CURDIR}/build

LAMBDA_PING_TEST_BINARY=LambdaPingTestHandler
LAMBDA_RELAY_TEST_BINARY=LambdaRelayTestHandler
LAMBDA_CORS_BINARY=LambdaCORSHandler
LAMBDA_PING_TEST_TARGET=./cmd/LambdaPingTestHandler
LAMBDA_CORS_TARGET=./cmd/LambdaCORSHandler
LAMBDA_RELAY_TEST_TARGET=./cmd/LambdaRelayTestHandler

PING_TEST_BINARY=pingtest
RELAY_TEST_BINARY=relaytest
PING_TEST_TARGET=./cmd/pingtest
RELAY_TEST_TARGET=./cmd/relaytest

AWS_CLI=docker run -v ~/.aws:/root/.aws -v ${PWD}/build:/build --rm -it amazon/aws-cli

help: ## Show this help message.
	echo "usage: make [target] ..."
	echo
	echo "targets:"
	echo "-------"
	egrep '^(.+)\:\ ##\ (.+)' Makefile | column -t -c 2 -s ':#'

docserver: ## Run an interactive OpenAPI spec on port 3333
	docker-compose up -d
	echo Visit documentation at http://localhost:3333

docserver-stop: ## Stop the interactive spec
	docker-compose down

docker-build-image: ## Builds the docker image
	docker build -t pokt-lint .

docker-remove-image: ## Removes the docker image
	docker rmi pokt-lint

docker-build-lambda: ## Builds the commands and Lambda functions for a linux environment
	make docker-build-lambda-pingtest
	make docker-build-lambda-relaytest
	echo ${BUILD_DIR}
	ls -l ${BUILD_DIR}

docker-build-lambda-pingtest: ## builds the pingtest lambda function
	docker run --rm -v ${BUILD_DIR}:/app/build pokt-lint go build -o ./build/${LAMBDA_PING_TEST_BINARY} ${LAMBDA_PING_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_PING_TEST_BINARY} && cd - > /dev/null

docker-build-lambda-relaytest: ## builds the relaytest lambda function
	docker run --rm -v ${BUILD_DIR}:/app/build pokt-lint go build -o ./build/${LAMBDA_RELAY_TEST_BINARY} ${LAMBDA_RELAY_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_RELAY_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY} && cd - > /dev/null

docker-build-commands: ## builds the commands for a linux environment
	make docker-build-pingtest
	make docker-build-relaytest

docker-build-pingtest: ## builds the pingtest command for a linux environment
	docker run --rm -v ${BUILD_DIR}:/app/build pokt-lint go build -o ./build/${PING_TEST_BINARY} ${PING_TEST_TARGET}

docker-build-relaytest: ## builds the relaytest command for a linux environment
	docker run --rm -v ${BUILD_DIR}:/app/build pokt-lint go build -o ./build/${RELAY_TEST_BINARY} ${RELAY_TEST_TARGET}

deploy-lambda-qa: ## builds and deploys the lambda QA functions
	make docker-build-lambda
	${AWS_CLI} lambda update-function-code --function-name qa-ping-test --zip-file fileb:///build/${LAMBDA_PING_TEST_BINARY}.zip
	${AWS_CLI} lambda update-function-code --function-name qa-relay-test --zip-file fileb:///build/${LAMBDA_RELAY_TEST_BINARY}.zip

deploy-lambda: ## builds and deploys the lambda function
	make docker-build-lambda
	${AWS_CLI} lambda update-function-code --function-name ping-test --zip-file fileb:///build/${LAMBDA_PING_TEST_BINARY}.zip
	${AWS_CLI} lambda update-function-code --function-name relay-test --zip-file fileb:///build/${LAMBDA_RELAY_TEST_BINARY}.zip

build-all: ## builds the commands and lambda functions for the local environment
	go build -o ./build/ ./cmd/...

test: ## runs the unit tests
	go test -v ./...

clean: ## deletes build artifacts
	go clean
	rm -f ${BUILD_DIR}/*

