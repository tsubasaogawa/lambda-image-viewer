package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func generateCamerarollHtml(key string) (events.LambdaFunctionURLResponse, error) {
	log.Printf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key)
	return responseHtml("", 200, Headers{"Cache-Control": "private"}), nil
}
