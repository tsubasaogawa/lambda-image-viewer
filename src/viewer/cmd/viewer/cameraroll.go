package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/url"
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
	ViewerDomain         string
	ImgWidthToClipboard  uint64
	ImgHeightToClipboard uint64
	LastKey              string
	PrevKeys             []string // Add this line to store the history of scan keys
	NextLink             string   // Add this line
	PrevLink             string   // Add this line
	IsPrivate            bool
	SaltForPrivateImage  string
}

func generateCamerarollHtml(currentScanKey dynamo.PagingKey, prevKeys []string, isPrivate bool) (events.LambdaFunctionURLResponse, error) {
	thumbs, lk, err := model.New().ListThumbnails(MAX_THUMBNAIL_PER_PAGE, currentScanKey)
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
	// Generate NextLink
	nextLink := ""
	if len(lk) > 0 {
		newPrevKeys := make([]string, len(prevKeys))
		copy(newPrevKeys, prevKeys)
		currentKeyStr := generateLastEvaluatedKeyQueryString(currentScanKey)
		if currentKeyStr != "" {
			newPrevKeys = append(newPrevKeys, currentKeyStr)
		}
		nextLink = fmt.Sprintf("/cameraroll/%s?prevKeys=%s", generateLastEvaluatedKeyQueryString(lk), url.QueryEscape(strings.Join(newPrevKeys, ",")))
		if isPrivate {
			nextLink += "&private=true"
		}
	}

	// Generate PrevLink
	prevLink := ""
	if len(prevKeys) > 0 {
		prevKeyToNavigate := prevKeys[len(prevKeys)-1]  // Get the last key from the history
		remainingPrevKeys := prevKeys[:len(prevKeys)-1] // Remove the last key from the history

		prevLink = fmt.Sprintf("/cameraroll/%s?prevKeys=%s", prevKeyToNavigate, url.QueryEscape(strings.Join(remainingPrevKeys, ",")))
		if isPrivate {
			prevLink += "&private=true"
		}
	} else if len(prevKeys) == 0 && len(currentScanKey) > 0 {
		prevLink = "/cameraroll/"
		if isPrivate {
			prevLink += "?private=true"
		}
	}

	if err = cr.Execute(w, CameraRollData{
		Thumbnails:           thumbs,
		OriginDomain:         os.Getenv("ORIGIN_DOMAIN"),
		ViewerDomain:         os.Getenv("VIEWER_DOMAIN"),
		ImgWidthToClipboard:  width,
		ImgHeightToClipboard: height,
		LastKey:              generateLastEvaluatedKeyQueryString(lk),
		PrevKeys:             prevKeys, // Pass the prevKeys slice to the template
		NextLink:             nextLink,
		PrevLink:             prevLink,
		IsPrivate:            isPrivate,
		SaltForPrivateImage:  os.Getenv("SALT_FOR_PRIVATE_IMAGE"),
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
