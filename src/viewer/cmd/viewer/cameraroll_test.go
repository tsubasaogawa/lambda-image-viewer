package main

import (
	"time"

	"github.com/guregu/dynamo"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

type MockDB struct{}

func (db *MockDB) ListThumbnails(max int64, scanKey dynamo.PagingKey) (*[]model.Thumbnail, dynamo.PagingKey, error) {
	timestamp := time.Now().Unix()
	thumbs := []model.Thumbnail{
		{Id: "test-id-1", Timestamp: timestamp},
		{Id: "test-id-2", Timestamp: timestamp + 1},
	}
	return &thumbs, nil, nil
}

func (db *MockDB) GetMetadata(id string) (*model.Metadata, error) {
	return &model.Metadata{
		Width:  100,
		Height: 100,
	}, nil
}
