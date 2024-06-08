package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

func generateMetadataJson(key string) (events.LambdaFunctionURLResponse, error) {
	meta, err := model.New().GetMetadata(key)
	if err != nil {
		msg := fmt.Sprintf(`{"error": "obtaining metadata error, but skips that", "details": "%s"}`, err.Error())
		log.Println(msg)
		return responseJson("{}", 200), nil
		// return responseJson(msg, 500), fmt.Errorf(msg)
	}

	_json, err := json.Marshal(*meta)
	if err != nil {
		msg := `{"error": "json marshal error"}`
		return responseJson(msg, 500), fmt.Errorf(msg)
	}

	return responseJson(string(_json), 200), nil
}
