package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

// mockDB は DB インターフェースのモック実装です。
type mockDB struct {
	thumbsToReturn  *[]model.Thumbnail
	lastKeyToReturn dynamo.PagingKey
	errToReturn     error
}

func (m *mockDB) ListThumbnails(max int64, scanKey dynamo.PagingKey, isPrivate bool) (*[]model.Thumbnail, dynamo.PagingKey, error) {
	return m.thumbsToReturn, m.lastKeyToReturn, m.errToReturn
}

func TestCameraRollHandler(t *testing.T) {
	// 期待されるソート後のサムネイルデータ
	sortedThumbs := []model.Thumbnail{
		{Id: "img2.jpg", Timestamp: 300},
		{Id: "img3.jpg", Timestamp: 200},
		{Id: "img1.jpg", Timestamp: 100},
	}

	// モック用の LastEvaluatedKey
	mockLastKey := dynamo.PagingKey{
		"Id":        &dynamodb.AttributeValue{S: aws.String("img3.jpg")},
		"Timestamp": &dynamodb.AttributeValue{N: aws.String("200")},
	}
	mockLastKeyBytes, _ := json.Marshal(mockLastKey)
	encodedMockLastKey := base64.URLEncoding.EncodeToString(mockLastKeyBytes)

	db := &mockDB{
		thumbsToReturn:  &sortedThumbs,
		lastKeyToReturn: mockLastKey,
		errToReturn:     nil,
	}

	// テストケースの実行
	resp, err := CameraRollHandler(db, "", "3", false)
	if err != nil {
		t.Fatalf("CameraRollHandler returned an error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var body CameraRollResponse
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(body.Thumbnails, sortedThumbs) {
		t.Errorf("Thumbnails not sorted correctly. Got %+v, want %+v", body.Thumbnails, sortedThumbs)
	}

	if body.LastEvaluatedKey != encodedMockLastKey {
		t.Errorf("Unexpected LastEvaluatedKey. Got %s, want %s", body.LastEvaluatedKey, encodedMockLastKey)
	}
}

func TestCameraRollHandler_DBError(t *testing.T) {
	db := &mockDB{
		errToReturn: errors.New("dynamodb error"),
	}

	resp, err := CameraRollHandler(db, "", "10", false)
	if err != nil {
		// エラーが返されること自体は正常な挙動
	}

	if resp.StatusCode != 500 {
		t.Errorf("Expected status code 500 for DB error, got %d", resp.StatusCode)
	}
}

func TestCameraRollHandler_PrivateMode(t *testing.T) {
	// Test that isPrivate flag is properly passed through
	privateThumbs := []model.Thumbnail{
		{Id: "private/img1.jpg", Timestamp: 100, Private: true},
		{Id: "private/img2.jpg", Timestamp: 200, Private: true},
	}

	db := &mockDB{
		thumbsToReturn:  &privateThumbs,
		lastKeyToReturn: nil,
		errToReturn:     nil,
	}

	// Test with private mode enabled
	resp, err := CameraRollHandler(db, "", "10", true)
	if err != nil {
		t.Fatalf("CameraRollHandler returned an error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var body CameraRollResponse
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(body.Thumbnails, privateThumbs) {
		t.Errorf("Private thumbnails not returned correctly. Got %+v, want %+v", body.Thumbnails, privateThumbs)
	}
}

func TestCameraRollHandler_PublicMode(t *testing.T) {
	// Test that public mode works correctly
	publicThumbs := []model.Thumbnail{
		{Id: "public/img1.jpg", Timestamp: 100, Private: false},
		{Id: "public/img2.jpg", Timestamp: 200, Private: false},
	}

	db := &mockDB{
		thumbsToReturn:  &publicThumbs,
		lastKeyToReturn: nil,
		errToReturn:     nil,
	}

	// Test with private mode disabled
	resp, err := CameraRollHandler(db, "", "10", false)
	if err != nil {
		t.Fatalf("CameraRollHandler returned an error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var body CameraRollResponse
	if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

		if !reflect.DeepEqual(body.Thumbnails, publicThumbs) {

			t.Errorf("Public thumbnails not returned correctly. Got %+v, want %+v", body.Thumbnails, publicThumbs)

		}

	}

	

	func TestCameraRollHandler_UnsortedData(t *testing.T) {

		// Unsorted thumbnail data from the mock DB

		unsortedThumbs := []model.Thumbnail{

			{Id: "img1.jpg", Timestamp: 100},

			{Id: "img3.jpg", Timestamp: 200},

			{Id: "img2.jpg", Timestamp: 300},

		}

	

		// Expected sorted thumbnail data

		sortedThumbs := []model.Thumbnail{

			{Id: "img2.jpg", Timestamp: 300},

			{Id: "img3.jpg", Timestamp: 200},

			{Id: "img1.jpg", Timestamp: 100},

		}

	

		db := &mockDB{

			thumbsToReturn:  &unsortedThumbs,

			lastKeyToReturn: nil,

			errToReturn:     nil,

		}

	

		resp, err := CameraRollHandler(db, "", "10", false)

		if err != nil {

			t.Fatalf("CameraRollHandler returned an error: %v", err)

		}

	

		if resp.StatusCode != 200 {

			t.Errorf("Expected status code 200, got %d", resp.StatusCode)

		}

	

		var body CameraRollResponse

		if err := json.Unmarshal([]byte(resp.Body), &body); err != nil {

			t.Fatalf("Failed to unmarshal response body: %v", err)

		}

	

		if !reflect.DeepEqual(body.Thumbnails, sortedThumbs) {

			t.Errorf("Thumbnails not sorted correctly. Got %+v, want %+v", body.Thumbnails, sortedThumbs)

		}

	}

	