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

type CameraRollData struct {
	Thumbnails   *[]model.Thumbnail
	OriginDomain string
}

func generateCamerarollHtml() (events.LambdaFunctionURLResponse, error) {
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
	if err = cr.Execute(w, CameraRollData{thumbs, os.Getenv("ORIGIN_DOMAIN")}); err != nil {
		return responseHtml("", 500, Headers{}), err
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private"}), nil
}
