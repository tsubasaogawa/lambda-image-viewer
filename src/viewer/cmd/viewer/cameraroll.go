package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const (
	DEFAULT_THUMBNAIL_PER_PAGE = 50
)

var (
	//go:embed templates/camera_roll.html.tmpl
	crTmpl string
)

type CameraRollData struct {
	Thumbnails          *[]model.Thumbnail
	OriginDomain        string
	ViewerDomain        string
	LastKey             string
	IsPrivate           bool
	SaltForPrivateImage string
}

type CameraRollResponse struct {
	Thumbnails       []model.Thumbnail `json:"thumbnails"`
	LastEvaluatedKey string            `json:"last_evaluated_key"`
}

type DB interface {
	ListThumbnails(max int64, scanKey dynamo.PagingKey, isPrivate bool) (*[]model.Thumbnail, dynamo.PagingKey, error)
}

// generateCameraRollHtml は、カメラロールの初期表示用のHTMLを生成します。
func generateCameraRollHtml(db DB, isPrivate bool) (events.LambdaFunctionURLResponse, error) {
	thumbs, lk, err := db.ListThumbnails(DEFAULT_THUMBNAIL_PER_PAGE, nil, isPrivate)
	if err != nil {
		log.Printf("Failed to list thumbnails for HTML: %v", err)
		return responseHtml("Internal server error", 500, Headers{}), err
	}

	var nextKey string
	if len(lk) > 0 {
		lkBytes, err := json.Marshal(lk)
		if err != nil {
			log.Printf("Failed to marshal last key for HTML: %v", err)
			return responseHtml("Internal server error", 500, Headers{}), err
		}
		nextKey = base64.URLEncoding.EncodeToString(lkBytes)
	}

	tmpl, err := template.New("cameraroll").Parse(crTmpl)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return responseHtml("Internal server error", 500, Headers{}), err
	}

	w := new(strings.Builder)
	data := CameraRollData{
		Thumbnails:          thumbs,
		OriginDomain:        os.Getenv("ORIGIN_DOMAIN"),
		ViewerDomain:        os.Getenv("VIEWER_DOMAIN"),
		LastKey:             nextKey,
		IsPrivate:           isPrivate,
		SaltForPrivateImage: os.Getenv("SALT_FOR_PRIVATE_IMAGE"),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return responseHtml("Internal server error", 500, Headers{}), err
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private"}), nil
}

// CameraRollHandler は、ページネーション用のJSON APIリクエストを処理します。
func CameraRollHandler(db DB, lastEvaluatedKey string, limitStr string, isPrivate bool) (events.LambdaFunctionURLResponse, error) {
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit == 0 {
		limit = DEFAULT_THUMBNAIL_PER_PAGE
	}

	var scanKey dynamo.PagingKey
	if lastEvaluatedKey != "" {
		decodedKey, err := base64.URLEncoding.DecodeString(lastEvaluatedKey)
		if err != nil {
			return responseJson("Invalid last_evaluated_key", 400, Headers{}), err
		}
		var keyMap map[string]*dynamodb.AttributeValue
		if err := json.Unmarshal(decodedKey, &keyMap); err != nil {
			return responseJson("Invalid last_evaluated_key format", 400, Headers{}), err
		}
		scanKey = keyMap
	}

	thumbs, lk, err := db.ListThumbnails(limit, scanKey, isPrivate)
	if err != nil {
		log.Printf("Failed to list thumbnails for API: %v", err)
		return responseJson("Internal server error", 500, Headers{}), err
	}

	var nextKey string
	if len(lk) > 0 {
		lkBytes, err := json.Marshal(lk)
		if err != nil {
			log.Printf("Failed to marshal last key for API: %v", err)
			return responseJson("Internal server error", 500, Headers{}), err
		}
		nextKey = base64.URLEncoding.EncodeToString(lkBytes)
	}

	var thumbnailsToShow []model.Thumbnail
	if thumbs != nil {
		thumbnailsToShow = *thumbs
	} else {
		thumbnailsToShow = []model.Thumbnail{}
	}

	res := CameraRollResponse{
		Thumbnails:       thumbnailsToShow,
		LastEvaluatedKey: nextKey,
	}

	resBody, err := json.Marshal(res)
	if err != nil {
		log.Printf("Failed to marshal response for API: %v", err)
		return responseJson("Internal server error", 500, Headers{}), err
	}

	return responseJson(string(resBody), 200, Headers{"Cache-Control": "private"}), nil
}
