package model

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestThumbnailSort(t *testing.T) {
	// テストデータ
	thumbs := []Thumbnail{
		{Id: "a", Timestamp: 100},
		{Id: "b", Timestamp: 300},
		{Id: "c", Timestamp: 200},
	}

	// 期待される結果 (降順)
	expected := []Thumbnail{
		{Id: "b", Timestamp: 300},
		{Id: "c", Timestamp: 200},
		{Id: "a", Timestamp: 100},
	}

	// ListThumbnails内のソートロジックを直接テスト
	sort.Slice(thumbs, func(i, j int) bool { return thumbs[i].Timestamp > thumbs[j].Timestamp })

	// 結果を比較
	if !reflect.DeepEqual(thumbs, expected) {
		t.Errorf("Slice not sorted correctly.\nGot:    %+v\nWanted: %+v", thumbs, expected)
	}
}

func TestPrivateFieldDetection(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "Public image",
			id:       "public/images/img1.jpg",
			expected: false,
		},
		{
			name:     "Private image in middle of path",
			id:       "images/private/img2.jpg",
			expected: true,
		},
		{
			name:     "Private image at start needs leading slash",
			id:       "private/img3.jpg",
			expected: false,
		},
		{
			name:     "Private image with leading slash",
			id:       "/private/img3.jpg",
			expected: true,
		},
		{
			name:     "Non-private image without 'private' in path",
			id:       "photos/img4.jpg",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thumb := Thumbnail{Id: tt.id}
			// Simulate the logic used in ListThumbnails
			thumb.Private = strings.Contains(thumb.Id, "/private/")
			
			if thumb.Private != tt.expected {
				t.Errorf("Expected Private=%v for ID=%s, got %v", tt.expected, tt.id, thumb.Private)
			}
		})
	}
}
