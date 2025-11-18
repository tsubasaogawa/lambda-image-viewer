package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/tagger"
)

func main() {
	lambda.Start(tagger.Index)
}
