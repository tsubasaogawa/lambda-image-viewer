package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

var getMetadata = func(key string) (*model.Metadata, error) {
	return model.New().GetMetadata(key)
}

func generateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	meta, err := getMetadata(key)
	if err != nil {
		msg := fmt.Sprintf(`{"error": "obtaining metadata error, but skips that", "details": "%s"}`, err.Error())
		log.Println(msg)
		headers := Headers{}
		if model.IsPrivateKey(key) {
			headers["Cache-Control"] = "private, no-store"
		}
		return responseJson("{}", 200, headers), nil
		// return responseJson(msg, 500), fmt.Errorf(msg)
	}

	_json, err := json.Marshal(*meta)
	if err != nil {
		msg := `{"error": "json marshal error"}`
		return responseJson(msg, 500, Headers{}), fmt.Errorf(msg)
	}

	cacheControl := "public, max-age=31536000"
	if model.IsPrivateKey(key) {
		cacheControl = "private, no-store"
	}
	return responseJson(string(_json), 200, Headers{"Cache-Control": cacheControl}), nil
}
