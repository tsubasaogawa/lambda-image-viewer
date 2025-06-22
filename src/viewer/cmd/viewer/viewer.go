package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"maps"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
)

type Headers map[string]string

var (
	imageGenerator      ImageGenerator      = &DefaultImageGenerator{}
	metadataGenerator   MetadataGenerator   = &DefaultMetadataGenerator{}
	camerarollGenerator CamerarollGenerator = &DefaultCamerarollGenerator{}
)

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
		return imageGenerator.GenerateImageHtml(key)
	case "metadata":
		return metadataGenerator.GenerateMetadataJson(key)
	case "cameraroll":
		q, err := url.ParseQuery(r.RawQueryString)
		if err != nil {
			msg := "query parsing error. query=" + r.RawQueryString
			return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
		}
		isPrivate := q.Get("private") == "true"

		var currentScanKey dynamo.PagingKey
		var prevKeys []string
		prevKeysParam := q.Get("prevKeys")
		if prevKeysParam != "" {
			var err error
			prevKeys, err = camerarollGenerator.DecompressPrevKeys(prevKeysParam)
			if err != nil {
				log.Printf("Failed to decompress prevKeys: %v", err)
				msg := fmt.Sprintf("failed to decompress prevKeys: %v", err)
				return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
			}
		}

		// Determine the currentScanKey based on the 'key' path parameter
		if key != "" {
			parts := strings.SplitN(key, "/", 2)
			id, _ := base64.URLEncoding.DecodeString(parts[0])
			ts, _ := base64.URLEncoding.DecodeString(parts[1])
			currentScanKey = dynamo.PagingKey{
				"Id":        &dynamodb.AttributeValue{S: aws.String(string(id))},
				"Timestamp": &dynamodb.AttributeValue{N: aws.String(string(ts))},
			}
		}

		return camerarollGenerator.GenerateCamerarollHtml(currentScanKey, prevKeys, isPrivate)
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
