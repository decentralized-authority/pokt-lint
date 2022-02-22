BUILD_DIR=./build
LAMBDA_BINARY=LambdaRequestHandler
PINGTEST_BINARY=pingtest
LAMBDA_TARGET=./cmd/LambdaRequestHandler
PINGTEST_TARGET=./cmd/pingtest

lambda:
	GOOS=linux go build -o ${BUILD_DIR}/${LAMBDA_BINARY} ${LAMBDA_TARGET}
	cd ${BUILD_DIR} && zip ${LAMBDA_BINARY}.zip ${LAMBDA_BINARY}

pingtest:
	go build -o ${BUILD_DIR}/${PINGTEST_BINARY} ${PINGTEST_TARGET}
	${BUILD_DIR}/${PINGTEST_BINARY}

clean:
	go clean
	rm -f ${LAMBDA_BINARY} ${PINGTEST_BINARY} ${LAMBDA_BINARY}.zip
