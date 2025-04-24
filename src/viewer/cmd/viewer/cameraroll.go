package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/guregu/dynamo"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const (
	MAX_THUMBNAIL_PER_PAGE = 200
)

var (
	//go:embed templates/camera_roll.html.tmpl
	crTmpl string
)

type CameraRollData struct {
	Thumbnails           *[]model.Thumbnail
	OriginDomain         string
	ImgWidthToClipboard  uint64
	ImgHeightToClipboard uint64
	LastKey              string
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
	width, _ := strconv.ParseUint(os.Getenv("IMG_WIDTH_TO_CLIPBOARD"), 10, 64)
	height, _ := strconv.ParseUint(os.Getenv("IMG_HEIGHT_TO_CLIPBOARD"), 10, 64)
	if err = cr.Execute(w, CameraRollData{
		thumbs,
		os.Getenv("ORIGIN_DOMAIN"),
		width,
		height,
		generateLastEvaluatedKeyQueryString(lk),
	}); err != nil {
		return responseHtml("", 500, Headers{}), err
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private"}), nil
}

func generateLastEvaluatedKeyQueryString(lk dynamo.PagingKey) string {
	if len(lk) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%s/%s",
		base64.URLEncoding.EncodeToString([]byte(*lk["Id"].S)),
		base64.URLEncoding.EncodeToString([]byte(*lk["Timestamp"].N)),
	)
}
