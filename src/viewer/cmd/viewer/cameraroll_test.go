package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
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

func TestGenerateCamerarollHtml(t *testing.T) {
	mockDB := &MockDB{}
	resp, err := generateCamerarollHtml(mockDB, nil, []string{}, false)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Get the timestamp from the mock
	thumbs, _, _ := mockDB.ListThumbnails(0, nil)
	expectedId := fmt.Sprintf(`id="%d"`, (*thumbs)[0].Timestamp)

	assert.True(t, strings.Contains(resp.Body, expectedId), "Expected HTML to contain id with timestamp")
}
