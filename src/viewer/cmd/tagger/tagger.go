package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Index)
}

func Index(ctx context.Context, event events.S3Event) {
	for _, r := range event.Records {
		fmt.Printf("s3://%s/%s\n", r.S3.Bucket.Name, r.S3.Object.Key)
	}
}
