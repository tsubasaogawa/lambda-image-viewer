package main

import (
	"fmt"
	"log"
	"maps"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

type Headers map[string]string

var (
	db                DB                = model.New()
	imageGenerator    ImageGenerator    = &DefaultImageGenerator{}
	metadataGenerator MetadataGenerator = &DefaultMetadataGenerator{}
)

func main() {
	lambda.Start(Index)
}

func Index(r events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	trimmedPath := strings.TrimPrefix(r.RawPath, "/")
	p := strings.SplitN(trimmedPath, "/", 2)

	// Handle cases like "/" or ""
	if len(p) < 1 || p[0] == "" {
		msg := "path parsing error. path=" + r.RawPath
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	route := p[0]
	var key string
	if len(p) > 1 {
		key = p[1]
	}
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
			return responseJson(msg, 500, Headers{}), fmt.Errorf(msg)
		}

		lastEvaluatedKey := q.Get("last_evaluated_key")
		limit := q.Get("limit")
		isPrivate := q.Get("private") == "true"

		// If last_evaluated_key is present, it's an API call for the next page.
		if lastEvaluatedKey != "" || limit != "" {
			return CameraRollHandler(db, lastEvaluatedKey, limit, isPrivate)
		}

		// Otherwise, it's a request for the initial HTML page.
		return generateCameraRollHtml(db, isPrivate)
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
