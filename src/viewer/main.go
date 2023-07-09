package main

import (
	_ "embed"

	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/models"
)

type Photo struct {
	Url      string
	Metadata models.Metadata
}

var (
	//go:embed templates/index.html.tmpl
	tmpl string
)

func main() {
	lambda.Start(index)
}

func index(r events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	key := strings.TrimPrefix(r.RawPath, "/")

	_tmpl, err := template.New("index").Parse(tmpl)
	if err != nil {
		msg := "template parsing error"
		return response(msg, 500), fmt.Errorf(msg)
	}

	meta, err := getMetadata(getId(key))
	if err != nil {
		msg := "obtaining metadata error"
		return response(msg, 500), fmt.Errorf(msg)
	}

	photo := Photo{
		Url:      "https://photo.ogatube.com/" + key,
		Metadata: meta,
	}

	w := new(strings.Builder)

	if err = _tmpl.Execute(w, photo); err != nil {
		msg := "template execution error"
		return response(msg, 500), fmt.Errorf(msg)
	}

	return response(w.String(), 200), nil
}

func response(body string, status int) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		Body: body,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=utf-8",
		},
		StatusCode: status,
	}
}

func getId(key string) string {
	extlen := len(filepath.Ext(key))
	return key[0 : len(key)-extlen]
}

func getMetadata(id string) (models.Metadata, error) {
	m, err := models.New().GetMetadata(id)
	return *m, err
}
