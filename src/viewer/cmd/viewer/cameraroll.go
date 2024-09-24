package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/guregu/dynamo"
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
	LastKey      string
}

func generateCamerarollHtml(scanKey dynamo.PagingKey) (events.LambdaFunctionURLResponse, error) {
	thumbs, lk, err := model.New().ListThumbnails(MAX_THUMBNAIL_PER_PAGE, scanKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lastEvaluatedKey = %+v\n", lk)

	cr, err := template.New("cameraroll").Parse(crTmpl)
	if err != nil {
		log.Fatal(err)
	}

	w := new(strings.Builder)
	if err = cr.Execute(w, CameraRollData{
		thumbs,
		os.Getenv("ORIGIN_DOMAIN"),
		fmt.Sprintf(
			"%s/%s",
			base64.URLEncoding.EncodeToString([]byte(*lk["Id"].S)),
			base64.URLEncoding.EncodeToString([]byte(*lk["Timestamp"].N)),
		),
	}); err != nil {
		return responseHtml("", 500, Headers{}), err
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private"}), nil
}
