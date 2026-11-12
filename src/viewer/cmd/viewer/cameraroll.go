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
	Thumbnails   []CameraRollThumbnail
	OriginDomain string
	ViewerDomain string
	LastKey      string
	IsPrivate    bool
}

type CameraRollThumbnail struct {
	model.Thumbnail
	ImageURL     string `json:"image_url,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

type CameraRollResponse struct {
	Thumbnails       []CameraRollThumbnail `json:"thumbnails"`
	LastEvaluatedKey string                `json:"last_evaluated_key"`
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
	thumbnails, err := prepareThumbnails(thumbs, isPrivate)
	if err != nil {
		log.Printf("Failed to sign thumbnails: %v", err)
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
		Thumbnails:   thumbnails,
		OriginDomain: os.Getenv("ORIGIN_DOMAIN"),
		ViewerDomain: os.Getenv("VIEWER_DOMAIN"),
		LastKey:      nextKey,
		IsPrivate:    isPrivate,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return responseHtml("Internal server error", 500, Headers{}), err
	}

	return responseHtml(w.String(), 200, Headers{"Cache-Control": "private, no-store"}), nil
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

	thumbnailsToShow, err := prepareThumbnails(thumbs, isPrivate)
	if err != nil {
		log.Printf("Failed to sign thumbnails: %v", err)
		return responseJson("Internal server error", 500, Headers{}), err
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

	return responseJson(string(resBody), 200, Headers{"Cache-Control": "private, no-store"}), nil
}

func prepareThumbnails(thumbs *[]model.Thumbnail, isPrivate bool) ([]CameraRollThumbnail, error) {
	result := []CameraRollThumbnail{}
	if thumbs == nil {
		return result, nil
	}
	var signer privateURLSigner
	for _, thumb := range *thumbs {
		thumb.Private = model.IsPrivateKey(thumb.Id)
		if thumb.Private != isPrivate {
			continue
		}
		entry := CameraRollThumbnail{Thumbnail: thumb}
		if thumb.Private {
			if signer == nil {
				var err error
				signer, err = newPrivateURLSigner()
				if err != nil {
					return nil, err
				}
			}
			urls, err := signPrivateURLs(signer, thumb.Id)
			if err != nil {
				return nil, err
			}
			entry.ImageURL, entry.ThumbnailURL = urls.ImageURL, urls.ThumbnailURL
		}
		result = append(result, entry)
	}
	return result, nil
}

func privateURLsHandler(key string) (events.LambdaFunctionURLResponse, error) {
	if !validPrivateKey(key) {
		return responseJson(`{"error":"invalid private image key"}`, 400, Headers{"Cache-Control": "private, no-store"}), nil
	}
	signer, err := newPrivateURLSigner()
	if err != nil {
		log.Printf("Failed to load signing key: %v", err)
		return responseJson(`{"error":"internal server error"}`, 500, Headers{"Cache-Control": "private, no-store"}), err
	}
	urls, err := signPrivateURLs(signer, key)
	if err != nil {
		log.Printf("Failed to sign private image URLs: %v", err)
		return responseJson(`{"error":"internal server error"}`, 500, Headers{"Cache-Control": "private, no-store"}), err
	}
	body, err := json.Marshal(urls)
	if err != nil {
		return responseJson(`{"error":"internal server error"}`, 500, Headers{"Cache-Control": "private, no-store"}), err
	}
	return responseJson(string(body), 200, Headers{"Cache-Control": "private, no-store"}), nil
}
