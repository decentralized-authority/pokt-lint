package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsnoproblem/pokt-lint/linting"
)

func main() {
	lambda.Start(linting.HandleRequest)
}
