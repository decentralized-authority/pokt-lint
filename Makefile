LAMBDA_BINARY=./build/LambdaRequestHandler
LAMBDA_PACKAGE=./cmd/LambdaRequestHandler
PINGTEST_BINARY=./build/pingtest
PINGTEST_PACKAGE=./cmd/pingtest

lambda:
	GOOS=linux go build -o ${LAMBDA_BINARY} ${LAMBDA_PACKAGE}
	zip ${LAMBDA_BINARY}.zip ${LAMBDA_BINARY}

pingtest:
	go build -o ${PINGTEST_BINARY} ${PINGTEST_PACKAGE}
	${PINGTEST_BINARY}

clean:
	go clean
	rm ${LAMBDA_BINARY} ${PINGTEST_BINARY} ${LAMBDA_BINARY}.zip
