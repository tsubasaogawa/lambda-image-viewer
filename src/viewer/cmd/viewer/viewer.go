package main

import (
	_ "embed"
	"encoding/json"

	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

type Image struct {
	Url      string
	Metadata model.Metadata
}

var (
	//go:embed templates/index.html.tmpl
	tmpl string
)

func main() {
	lambda.Start(Index)
}

func Index(r events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	p := strings.SplitN(strings.TrimPrefix(r.RawPath, "/"), "/", 2)
	if p == nil {
		msg := "path parsing error"
		return responseHtml(msg, 500), fmt.Errorf(msg)
	}
	route := p[0]
	key := p[1]

	switch route {
	case "image":
		return generateImageHtml(key)
	case "metadata":
		return generateMetadataJson(key)
	default:
		msg := "no route error"
		return responseHtml(msg, 500), fmt.Errorf(msg)
	}
}

func generateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	_tmpl, err := template.New("index").Parse(tmpl)
	if err != nil {
		msg := "template parsing error"
		return responseHtml(msg, 500), fmt.Errorf(msg)
	}

	meta, err := model.New().GetMetadata(getId(key))
	if err != nil {
		msg := "obtaining metadata error"
		return responseHtml(msg, 500), fmt.Errorf(msg)
	}

	image := Image{
		Url:      fmt.Sprintf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key),
		Metadata: *meta,
	}

	w := new(strings.Builder)

	if err = _tmpl.Execute(w, image); err != nil {
		msg := "template execution error"
		return responseHtml(msg, 500), fmt.Errorf(msg)
	}

	return responseHtml(w.String(), 200), nil
}

func generateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	meta, err := model.New().GetMetadata(getId(key))
	if err != nil {
		msg := `{"error": "obtaining metadata error"}`
		return responseJson(msg, 500), fmt.Errorf(msg)
	}

	_json, err := json.Marshal(*meta)
	if err != nil {
		msg := `{"error": "json marshal error"}`
		return responseJson(msg, 500), fmt.Errorf(msg)
	}

	return responseJson(string(_json), 200), nil
}

func response(body, _type string, status int) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body: body,
		Headers: map[string]string{
			"Content-Type": _type,
		},
		StatusCode: status,
	}
}

func responseHtml(body string, status int) events.LambdaFunctionURLResponse {
	return response(body, "text/html; charset=utf-8", status)
}

func responseJson(body string, status int) events.LambdaFunctionURLResponse {
	return response(body, "application/json; charset=utf-8", status)
}

func getId(key string) string {
	return key
}
