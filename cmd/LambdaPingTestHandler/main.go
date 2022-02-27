package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsnoproblem/pokt-lint/pinging"
)

func main() {
	lambda.Start(pinging.HandleRequest)
}
