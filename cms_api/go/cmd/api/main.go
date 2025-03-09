package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
		Name string `json:"name"`
}
 
func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, name MyEvent) (string, error) {
    return fmt.Sprintf("Hello %s!", name.Name), nil
}
 