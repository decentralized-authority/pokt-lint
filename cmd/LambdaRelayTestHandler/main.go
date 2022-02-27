package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsnoproblem/pokt-lint/relaying"
)

func main() {
	lambda.Start(relaying.HandleRequest)
}
