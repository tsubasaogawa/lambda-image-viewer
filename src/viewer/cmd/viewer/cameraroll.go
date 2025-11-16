package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const (
	DEFAULT_THUMBNAIL_PER_PAGE = 50
)

type CameraRollResponse struct {
	Thumbnails       []model.Thumbnail `json:"thumbnails"`
	LastEvaluatedKey string            `json:"last_evaluated_key"`
}

type DB interface {
	ListThumbnails(max int64, scanKey dynamo.PagingKey) (*[]model.Thumbnail, dynamo.PagingKey, error)
}

func CameraRollHandler(db DB, lastEvaluatedKey string, limitStr string) (events.LambdaFunctionURLResponse, error) {
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

	thumbs, lk, err := db.ListThumbnails(limit, scanKey)
	if err != nil {
		log.Printf("Failed to list thumbnails: %v", err)
		return responseJson("Internal server error", 500, Headers{}), err
	}

	var nextKey string
	if len(lk) > 0 {
		lkBytes, err := json.Marshal(lk)
		if err != nil {
			log.Printf("Failed to marshal last key: %v", err)
			return responseJson("Internal server error", 500, Headers{}), err
		}
		nextKey = base64.URLEncoding.EncodeToString(lkBytes)
	}

	// `thumbs`がnilの場合に空のスライスを割り当てる
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
		log.Printf("Failed to marshal response: %v", err)
		return responseJson("Internal server error", 500, Headers{}), err
	}

	return responseJson(string(resBody), 200, Headers{"Cache-Control": "private"}), nil
}
