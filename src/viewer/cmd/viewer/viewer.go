package main

import (
	"maps"

	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Headers map[string]string

func main() {
	lambda.Start(Index)
}

func Index(r events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	p := strings.SplitN(strings.TrimPrefix(r.RawPath, "/"), "/", 2)
	if p == nil || len(p) < 2 {
		msg := "path parsing error. path=" + r.RawPath
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}
	route := p[0]
	key := p[1]

	switch route {
	case "image":
		return generateImageHtml(key)
	case "metadata":
		return generateMetadataJson(key)
	case "cameraroll":
		return generateCamerarollHtml(key)
	default:
		msg := "no route error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}
}

func response(body string, status int, headers Headers) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body:       body,
		Headers:    headers,
		StatusCode: status,
	}
}

func responseHtml(body string, status int, headers Headers) events.LambdaFunctionURLResponse {
	maps.Copy(headers, Headers{"Content-Type": "text/html; charset=utf-8"})
	return response(body, status, headers)
}

func responseJson(body string, status int) events.LambdaFunctionURLResponse {
	return response(body, status, Headers{"Content-Type": "application/json; charset=utf-8"})
}
