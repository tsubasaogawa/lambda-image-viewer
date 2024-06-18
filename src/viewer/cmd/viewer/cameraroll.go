package main

import (
	_ "embed"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const (
	MAX_THUMBNAIL_PER_PAGE = 20
)

var (
	//go:embed templates/camera_roll.html.tmpl
	crTmpl string
)

func generateCamerarollHtml(key string) (events.LambdaFunctionURLResponse, error) {
	log.Printf("https://%s/%s", os.Getenv("ORIGIN_DOMAIN"), key)
	thumbs, lk, err := model.New().ListThumbnails(MAX_THUMBNAIL_PER_PAGE, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lastEvaluatedKey = %+v\n", lk)

	cr, err := template.New("cameraroll").Parse(crTmpl)
	if err != nil {
		log.Fatal(err)
	}

	w := new(strings.Builder)
	cr.Execute(w, thumbs)

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private"}), nil
}
