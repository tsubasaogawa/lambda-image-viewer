package main

import (
	_ "embed"

	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

func generateImageHtml(key string) (events.LambdaFunctionURLResponse, error) {
	_tmpl, err := template.New("index").Parse(tmpl)
	if err != nil {
		msg := "template parsing error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	meta, err := model.New().GetMetadata(key)
	if err != nil {
		log.Println("obtaining metadata error. viewer uses empty metadata.")
		meta = &model.Metadata{}
	}

	image := Image{
		Url:      fmt.Sprintf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key),
		Metadata: *meta,
	}

	w := new(strings.Builder)

	if err = _tmpl.Execute(w, image); err != nil {
		msg := "template execution error"
		return responseHtml(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	return responseHtml(w.String(), 200, Headers{}), nil
}
