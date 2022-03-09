package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
)

func main() {
	lambda.Start(handleCORSRequest)
}

func handleCORSRequest(ctx context.Context, _ interface{}) (http.Response, error) {
	header := http.Header{}
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	header.Set("Access-Control-Allow-Headers", "Content-Type")
	body := ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	resp := http.Response{
		Status:        "OK",
		StatusCode:    200,
		Header:        header,
		Body:          body,
		ContentLength: 0,
	}
	return resp, nil
}
