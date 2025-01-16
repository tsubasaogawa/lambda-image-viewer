package main

import (
	"encoding/base64"
	"log"
	"maps"

	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
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
	log.Printf("route=%s, key=%s\n", route, key)

	switch route {
	case "image":
		return generateImageHtml(key)
	case "metadata":
		return generateMetadataJson(key)
	case "cameraroll":
		if !authorized(r.Headers) {
			h := Headers{}
			h["WWW-Authenticate"] = "Basic"
			return responseHtml("Unauthorized", 401, h), fmt.Errorf("Unauthorized")
		}

		if key == "" {
			return generateCamerarollHtml(nil)
		}
		parts := strings.SplitN(key, "/", 2)
		id, _ := base64.URLEncoding.DecodeString(parts[0])
		ts, _ := base64.URLEncoding.DecodeString(parts[1])

		return generateCamerarollHtml(dynamo.PagingKey{
			"Id":        &dynamodb.AttributeValue{S: aws.String(string(id))},
			"Timestamp": &dynamodb.AttributeValue{N: aws.String(string(ts))},
		})
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

func responseJson(body string, status int, headers Headers) events.LambdaFunctionURLResponse {
	maps.Copy(headers, Headers{"Content-Type": "application/json; charset=utf-8"})
	return response(body, status, headers)
}
