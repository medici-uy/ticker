package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func handler() error {
	return nil
}

func main() {
	lambda.Start(handler)
}
