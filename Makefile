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

docker-build-commands: ## Builds the commands
	docker run --rm -ti -v ${BUILD_DIR}:/app/build pokt-lint build-commands

docker-build-lambda: ## Builds the lambda functions
	docker run --rm -ti -v ${BUILD_DIR}:/app/build pokt-lint build-lambda

build-commands: ## <-- compiles executables to ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${PING_TEST_BINARY} ${PING_TEST_TARGET}
	go build -o ${BUILD_DIR}/${RELAY_TEST_BINARY} ${RELAY_TEST_TARGET}

build-lambda: ## <-- builds lambda function bundles in ${BUILD_DIR}
	make build-lambda-pingtest
	make build-lambda-relaytest
	make build-lambda-cors

build-lambda-pingtest: ## builds the pingtest lambda function
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_PING_TEST_BINARY} ${LAMBDA_PING_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_PING_TEST_BINARY}

build-lambda-relaytest: ## builds the relaytest lambda function
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_RELAY_TEST_BINARY} ${LAMBDA_RELAY_TEST_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_RELAY_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY}

build-lambda-cors: ## builds the cors handler (to return access-control-* headers)
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_CORS_BINARY} ${LAMBDA_CORS_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_CORS_BINARY}.zip ${LAMBDA_CORS_BINARY}

deploy-lambda-all-qa: ## <-- deploys all lambda QA functions
	make deploy-lambda-cors-qa
	make deploy-lambda-pingtest-qa
	make deploy-lambda-relaytest-qa

deploy-lambda-all: ## <-- deploys all lambda functions
	make deploy-lambda-cors
	make deploy-lambda-pingtest
	make deploy-lambda-relaytest

deploy-lambda-pingtest-qa: ## builds and deploys the pingtest QA function
	make build-lambda-pingtest
	${AWS_CLI} lambda update-function-code --function-name qa-ping-test --zip-file fileb:///build/${LAMBDA_PING_TEST_BINARY}.zip

deploy-lambda-relaytest-qa: ## builds and deploys the relaytest QA function
	make build-lambda-relaytest
	${AWS_CLI} lambda update-function-code --function-name qa-relay-test --zip-file fileb:///build/${LAMBDA_RELAY_TEST_BINARY}.zip

deploy-lambda-cors-qa: ## builds and deploys the CORS QA function
	make build-lambda-cors
	${AWS_CLI} lambda update-function-code --function-name qa-cors-handler --zip-file fileb:///build/${LAMBDA_CORS_BINARY}.zip

deploy-lambda-pingtest: ## builds and deploys the pingtest function
	make build-lambda-pingtest
	${AWS_CLI} lambda update-function-code --function-name ping-test --zip-file fileb:///build/${LAMBDA_PING_TEST_BINARY}.zip

deploy-lambda-relaytest: ## builds and deploys the relaytest function
	make build-lambda-relaytest
	${AWS_CLI} lambda update-function-code --function-name relay-test --zip-file fileb:///build/${LAMBDA_RELAY_TEST_BINARY}.zip

deploy-lambda-cors: ## builds and deploys the CORS function
	make build-lambda-cors
	${AWS_CLI} lambda update-function-code --function-name cors-handler --zip-file fileb:///build/${LAMBDA_CORS_BINARY}.zip

test: ## runs the unit tests
	go test -v ./...

clean: ## deletes build artifacts
	go clean
	rm -f ${LAMBDA_PING_TEST_BINARY} ${LAMBDA_RELAY_TEST_BINARY} ${PING_TEST_BINARY} ${RELAY_TEST_BINARY} ${LAMBDA_PING_TEST_BINARY}.zip ${LAMBDA_RELAY_TEST_BINARY}.zip
	rm -f ${BUILD_DIR}/${LAMBDA_CORS_BINARY} ${BUILD_DIR}/${LAMBDA_CORS_BINARY}.zip
